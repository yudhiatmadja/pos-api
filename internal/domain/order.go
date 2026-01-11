package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusNew      OrderStatus = "NEW"
	OrderStatusAccepted OrderStatus = "ACCEPTED"
	OrderStatusCooking  OrderStatus = "COOKING"
	OrderStatusReady    OrderStatus = "READY"
	OrderStatusDone     OrderStatus = "DONE"
	OrderStatusVoided   OrderStatus = "VOIDED"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid   PaymentStatus = "UNPAID"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

type Order struct {
	ID             uuid.UUID     `json:"id"`
	StoreID        uuid.UUID     `json:"store_id"`
	TableSessionID *uuid.UUID    `json:"table_session_id,omitempty"`
	CashierID      *uuid.UUID    `json:"cashier_id,omitempty"`
	OrderNumber    string        `json:"order_number"`
	Status         OrderStatus   `json:"status"`
	PaymentStatus  PaymentStatus `json:"payment_status"`
	TotalAmount    float64       `json:"total_amount"`
	TaxAmount      float64       `json:"tax_amount"`
	DiscountAmount float64       `json:"discount_amount"`
	FinalAmount    float64       `json:"final_amount"`
	Note           string        `json:"note,omitempty"`
	Items          []OrderItem   `json:"items"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type OrderItem struct {
	ID           uuid.UUID `json:"id"`
	OrderID      uuid.UUID `json:"order_id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice float64   `json:"product_price"`
	Quantity     int32     `json:"quantity"`
	TotalPrice   float64   `json:"total_price"`
	Note         string    `json:"note,omitempty"`
}

type CreateOrderRequest struct {
	StoreID        uuid.UUID                `json:"store_id" binding:"required"`
	TableSessionID *uuid.UUID               `json:"table_session_id"`
	Note           string                   `json:"note"`
	Items          []CreateOrderItemRequest `json:"items" binding:"required,dive"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int32     `json:"quantity" binding:"required,gt=0"`
	Note      string    `json:"note"`
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*Order, error)
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status OrderStatus, userID uuid.UUID) (*Order, error)
	GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]Order, error)
}
