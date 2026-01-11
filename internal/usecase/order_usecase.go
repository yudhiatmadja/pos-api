package usecase

import (
	"context" // Keeping context as it's used throughout the file. The instruction to remove it seems to be based on a misunderstanding or an incomplete example.
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"pos-api/internal/domain"
	"pos-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderUsecase struct {
	store    repository.Repository
	eventSvc domain.EventService
}

func NewOrderUsecase(store repository.Repository, eventSvc domain.EventService) domain.OrderUsecase {
	return &orderUsecase{
		store:    store,
		eventSvc: eventSvc,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error) {
	// 1. Idempotency Check
	if req.IdempotencyKey != "" {
		existing, err := uc.store.GetIdempotencyKey(ctx, req.IdempotencyKey)
		if err == nil {
			// Found existing key, return cached order
			var cachedOrder domain.Order
			if jsonErr := json.Unmarshal(existing.ResponseBody, &cachedOrder); jsonErr == nil {
				return &cachedOrder, nil
			}
			// If unmarshal fails, we proceed to recreate (or log error) - treating as new for safety or erroring?
			// Safer to error to warn client.
			return nil, fmt.Errorf("failed to recover idempotent order")
		}
	}

	var order domain.Order

	err := uc.store.ExecTx(ctx, func(q *repository.Queries) error {
		// 1. Validate Products & Calculate Total
		var totalAmount float64
		var orderItems []domain.OrderItem

		for _, itemReq := range req.Items {
			product, err := q.GetProduct(ctx, pgtype.UUID{Bytes: itemReq.ProductID, Valid: true})
			if err != nil {
				return fmt.Errorf("product not found: %s", itemReq.ProductID)
			}

			if !product.IsAvailable.Bool || product.Stock < itemReq.Quantity {
				return fmt.Errorf("product not available or insufficient stock: %s", product.Name)
			}

			price, _ := product.Price.Float64Value()
			itemTotal := price.Float64 * float64(itemReq.Quantity)
			totalAmount += itemTotal

			orderItems = append(orderItems, domain.OrderItem{
				ProductID:    uuid.UUID(product.ID.Bytes),
				ProductName:  product.Name,
				ProductPrice: price.Float64,
				Quantity:     itemReq.Quantity,
				TotalPrice:   itemTotal,
				Note:         itemReq.Note,
			})
		}

		// 2. Create Order Header
		totalAmountNumeric := pgtype.Numeric{Int: big.NewInt(int64(totalAmount * 100)), Exp: -2, Valid: true}
		orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())

		var sessionID pgtype.UUID
		if req.TableSessionID != nil {
			sessionID = pgtype.UUID{Bytes: *req.TableSessionID, Valid: true}
		}

		// Use StoreID instead of OutletID
		dbOrder, err := q.CreateOrder(ctx, repository.CreateOrderParams{
			StoreID:        pgtype.UUID{Bytes: req.StoreID, Valid: true},
			TableSessionID: sessionID,
			OrderNumber:    orderNumber,
			TotalAmount:    totalAmountNumeric,
			FinalAmount:    totalAmountNumeric,
			Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
			// Status and PaymentStatus default in DB
		})
		if err != nil {
			return err
		}

		// 3. Create Order Items
		for _, item := range orderItems {
			itemPrice := pgtype.Numeric{Int: big.NewInt(int64(item.ProductPrice * 100)), Exp: -2, Valid: true}
			itemTotal := pgtype.Numeric{Int: big.NewInt(int64(item.TotalPrice * 100)), Exp: -2, Valid: true}

			_, err := q.CreateOrderItem(ctx, repository.CreateOrderItemParams{
				OrderID:      dbOrder.ID,
				ProductID:    pgtype.UUID{Bytes: item.ProductID, Valid: true},
				ProductName:  item.ProductName,
				ProductPrice: itemPrice,
				Quantity:     item.Quantity,
				TotalPrice:   itemTotal,
				Note:         pgtype.Text{String: item.Note, Valid: item.Note != ""},
			})
			if err != nil {
				return err
			}
		}

		// Populate return struct
		tVal, _ := dbOrder.TotalAmount.Float64Value()
		order = domain.Order{
			ID:          uuid.UUID(dbOrder.ID.Bytes),
			StoreID:     uuid.UUID(dbOrder.StoreID.Bytes),
			OrderNumber: dbOrder.OrderNumber,
			Status:      domain.OrderStatus(dbOrder.Status),
			TotalAmount: tVal.Float64,
			Items:       orderItems,
			CreatedAt:   dbOrder.CreatedAt.Time,
		}

		// 4. Save Idempotency Key (Inside Tx for consistency)
		if req.IdempotencyKey != "" {
			jsonBytes, _ := json.Marshal(order)
			_, err = q.CreateIdempotencyKey(ctx, repository.CreateIdempotencyKeyParams{
				Key:            req.IdempotencyKey,
				ResponseStatus: 201,
				ResponseBody:   jsonBytes,
			})
			if err != nil {
				return fmt.Errorf("failed to save idempotency key: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish Realtime Event
	_ = uc.eventSvc.PublishEvent(ctx, "NEW_ORDER", order)

	return &order, nil
}

func (uc *orderUsecase) UpdateStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus, userID uuid.UUID) (*domain.Order, error) {
	// 1. Get current order to validate transition
	currentOrder, err := uc.store.GetOrder(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("order not found")
	}

	// 2. Validate state transition
	if !isValidTransition(domain.OrderStatus(currentOrder.Status), status) {
		return nil, fmt.Errorf("invalid status transition from %s to %s", currentOrder.Status, status)
	}

	// 3. Update status
	dbOrder, err := uc.store.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
		ID:     pgtype.UUID{Bytes: orderID, Valid: true},
		Status: string(status),
	})
	if err != nil {
		return nil, err
	}

	// 4. Audit Log
	_, _ = uc.store.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		UserID:   pgtype.UUID{Bytes: userID, Valid: true},
		Action:   "UPDATE_ORDER_STATUS",
		Entity:   pgtype.Text{String: "Order", Valid: true},
		EntityID: pgtype.UUID{Bytes: orderID, Valid: true},
		Before:   []byte(fmt.Sprintf(`{"status": "%s"}`, currentOrder.Status)),
		After:    []byte(fmt.Sprintf(`{"status": "%s"}`, status)),
	})

	// Create struct
	return &domain.Order{
		ID:     uuid.UUID(dbOrder.ID.Bytes),
		Status: domain.OrderStatus(dbOrder.Status),
	}, nil
}

func (uc *orderUsecase) GetOrder(ctx context.Context, orderID uuid.UUID) (*domain.Order, error) {
	dbOrder, err := uc.store.GetOrder(ctx, pgtype.UUID{Bytes: orderID, Valid: true})
	if err != nil {
		return nil, err
	}
	// Fetch items too... kept simple for now
	return &domain.Order{
		ID:     uuid.UUID(dbOrder.ID.Bytes),
		Status: domain.OrderStatus(dbOrder.Status),
	}, nil
}

func (uc *orderUsecase) GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]domain.Order, error) {
	orders, err := uc.store.GetOrdersBySession(ctx, pgtype.UUID{Bytes: sessionID, Valid: true})
	if err != nil {
		return nil, err
	}

	var res []domain.Order
	for _, o := range orders {
		res = append(res, domain.Order{
			ID:          uuid.UUID(o.ID.Bytes),
			Status:      domain.OrderStatus(o.Status),
			OrderNumber: o.OrderNumber,
		})
	}
	return res, nil
}
