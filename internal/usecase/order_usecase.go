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
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderUsecase struct {
	repo *repository.Queries
	pool *pgxpool.Pool
}

func NewOrderUsecase(repo *repository.Queries, pool *pgxpool.Pool) domain.OrderUsecase {
	return &orderUsecase{
		repo: repo,
		pool: pool,
	}
}

func (uc *orderUsecase) CreateOrder(ctx context.Context, req *domain.CreateOrderRequest) (*domain.Order, error) {
	// Start transaction
	tx, err := uc.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	qtx := uc.repo.WithTx(tx)

	// 1. Validate Products & Calculate Total
	var totalAmount float64
	var orderItems []domain.OrderItem

	for _, itemReq := range req.Items {
		product, err := qtx.GetProduct(ctx, pgtype.UUID{Bytes: itemReq.ProductID, Valid: true})
		if err != nil {
			return nil, fmt.Errorf("product not found: %s", itemReq.ProductID)
		}

		if !product.IsAvailable.Bool || product.Stock < itemReq.Quantity {
			return nil, fmt.Errorf("product not available or insufficient stock: %s", product.Name)
		}

		price, _ := product.Price.Float64Value() // Numeric to Float64
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

	// 2. Create Order
	// Numeric handling in pgx can be tricky, simplifying for this example
	// In production, use precise decimal types.

	totalAmountNumeric := pgtype.Numeric{Int: big.NewInt(int64(totalAmount * 100)), Exp: -2, Valid: true}

	// Generate Order Number (Simple timestamp based or UUID)
	orderNumber := fmt.Sprintf("ORD-%d", time.Now().Unix())

	var sessionID pgtype.UUID
	if req.TableSessionID != nil {
		sessionID = pgtype.UUID{Bytes: *req.TableSessionID, Valid: true}
	}

	order, err := qtx.CreateOrder(ctx, repository.CreateOrderParams{
		OutletID:       pgtype.UUID{Bytes: req.OutletID, Valid: true},
		TableSessionID: sessionID,
		OrderNumber:    orderNumber,
		TotalAmount:    totalAmountNumeric,
		FinalAmount:    totalAmountNumeric, // No tax/discount logic yet
		Note:           pgtype.Text{String: req.Note, Valid: req.Note != ""},
	})
	if err != nil {
		return nil, err
	}

	// 3. Create Order Items
	for _, item := range orderItems {
		itemPrice := pgtype.Numeric{Int: big.NewInt(int64(item.ProductPrice * 100)), Exp: -2, Valid: true}
		itemTotal := pgtype.Numeric{Int: big.NewInt(int64(item.TotalPrice * 100)), Exp: -2, Valid: true}

		_, err := qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    pgtype.UUID{Bytes: item.ProductID, Valid: true},
			ProductName:  item.ProductName,
			ProductPrice: itemPrice,
			Quantity:     item.Quantity,
			TotalPrice:   itemTotal,
			Note:         pgtype.Text{String: item.Note, Valid: item.Note != ""},
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	// TODO: Publish Event to Redis

	totalVal, _ := order.TotalAmount.Float64Value()

	return &domain.Order{
		ID:          uuid.UUID(order.ID.Bytes),
		OutletID:    uuid.UUID(order.OutletID.Bytes),
		OrderNumber: order.OrderNumber,
		Status:      domain.OrderStatus(order.Status),
		TotalAmount: totalVal.Float64,
		Items:       orderItems,
	}, nil
}
