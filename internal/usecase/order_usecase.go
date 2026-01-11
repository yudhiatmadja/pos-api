package usecase

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"pos-api/internal/domain"
	"pos-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type orderUsecase struct {
	store    repository.Store
	eventSvc domain.EventService
}

func NewOrderUsecase(store repository.Store, eventSvc domain.EventService) domain.OrderUsecase {
	return &orderUsecase{
		store:    store,
		eventSvc: eventSvc,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error) {
	// TODO: Idempotency Check

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

			// Optimistic locking for stock could be adding UpdateProductStock query
			// For now assuming stock check is sufficient or ignoring race condition for MVP

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

		dbOrder, err := q.CreateOrder(ctx, repository.CreateOrderParams{
			OutletID:       pgtype.UUID{Bytes: req.OutletID, Valid: true},
			TableSessionID: sessionID,
			OrderNumber:    orderNumber,
			TotalAmount:    totalAmountNumeric,
			FinalAmount:    totalAmountNumeric,
			Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
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
			OutletID:    uuid.UUID(dbOrder.OutletID.Bytes),
			OrderNumber: dbOrder.OrderNumber,
			Status:      domain.OrderStatus(dbOrder.Status),
			TotalAmount: tVal.Float64,
			Items:       orderItems,
			CreatedAt:   dbOrder.CreatedAt.Time,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Publish Realtime Event
	_ = uc.eventSvc.PublishEvent(ctx, "NEW_ORDER", order)

	// TODO: Publish Realtime Event (NEW_ORDER)
	return &order, nil
}

func (uc *orderUsecase) UpdateStatus(ctx context.Context, orderID uuid.UUID, status domain.OrderStatus, userID uuid.UUID) (*domain.Order, error) {
	// 1. Get current order to validate transition
	// 2. Validate state transition
	// 3. Update status
	// 4. Publish event

	// Simplified implementation
	dbOrder, err := uc.store.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
		ID:     pgtype.UUID{Bytes: orderID, Valid: true},
		Status: string(status),
	})
	if err != nil {
		return nil, err // Handle error better
	}

	_ = userID // Audit log integration needed here

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
