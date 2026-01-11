package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string
type PaymentStatus string

const (
	OrderStatusNew      OrderStatus = "NEW"
	OrderStatusAccepted OrderStatus = "ACCEPTED"
	OrderStatusCooking  OrderStatus = "COOKING"
	OrderStatusReady    OrderStatus = "READY"
	OrderStatusDone     OrderStatus = "DONE"
	OrderStatusVoided   OrderStatus = "VOIDED"
)

const (
	PaymentStatusUnpaid   PaymentStatus = "UNPAID"
	PaymentStatusPaid     PaymentStatus = "PAID"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

type Order struct {
	ID             uuid.UUID     `json:"id"`
	OutletID       uuid.UUID     `json:"outlet_id"`
	TableSessionID *uuid.UUID    `json:"table_session_id,omitempty"`
	CashierID      *uuid.UUID    `json:"cashier_id,omitempty"`
	OrderNumber    string        `json:"order_number"`
	Status         OrderStatus   `json:"status"`
	PaymentStatus  PaymentStatus `json:"payment_status"`
	TotalAmount    float64       `json:"total_amount"`
	FinalAmount    float64       `json:"final_amount"`
	Note           string        `json:"note"`
	Items          []OrderItem   `json:"items,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
}

type OrderItem struct {
	ID           uuid.UUID `json:"id"`
	OrderID      uuid.UUID `json:"order_id"`
	ProductID    uuid.UUID `json:"product_id"`
	ProductName  string    `json:"product_name"`
	ProductPrice float64   `json:"product_price"`
	Quantity     int32     `json:"quantity"`
	TotalPrice   float64   `json:"total_price"`
	Note         string    `json:"note"`
}

type CreateOrderRequest struct {
	OutletID       uuid.UUID                `json:"outlet_id" binding:"required"`
	Items          []CreateOrderItemRequest `json:"items" binding:"required,min=1"`
	TableSessionID *uuid.UUID               `json:"table_session_id"` // Optional (if QR)
	Note           string                   `json:"note"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int32     `json:"quantity" binding:"required,min=1"`
	Note      string    `json:"note"`
}

type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*Order, error)
	// Order Processing
	UpdateStatus(ctx context.Context, orderID uuid.UUID, status OrderStatus, userID uuid.UUID) (*Order, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*Order, error)
	GetOrdersBySession(ctx context.Context, sessionID uuid.UUID) ([]Order, error)
}
