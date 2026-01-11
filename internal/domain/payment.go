package domain

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	PaymentMethodCash     PaymentMethod = "CASH"
	PaymentMethodQRIS     PaymentMethod = "QRIS"
	PaymentMethodPayLater PaymentMethod = "PAY_LATER"
)

type PaymentStatusType string // Rename from PaymentStatus to avoid conflict if in same package as Order (unlikely conflict if carefully named, but safe)
const (
	PaymentPending PaymentStatusType = "PENDING"
	PaymentSuccess PaymentStatusType = "SUCCESS"
	PaymentFailed  PaymentStatusType = "FAILED"
)

type Payment struct {
	ID              uuid.UUID         `json:"id"`
	OrderID         uuid.UUID         `json:"order_id"`
	PaymentMethod   PaymentMethod     `json:"payment_method"`
	Amount          float64           `json:"amount"`
	ReferenceNumber string            `json:"reference_number,omitempty"`
	QRISImageURL    string            `json:"qris_image_url,omitempty"`
	Status          PaymentStatusType `json:"status"`
	PaidAt          *time.Time        `json:"paid_at,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

type UploadQRISRequest struct {
	OrderID uuid.UUID             `form:"order_id" binding:"required"`
	File    *multipart.FileHeader `form:"image" binding:"required"`
}

type PaymentUsecase interface {
	UploadQRIS(ctx context.Context, req *UploadQRISRequest) (*Payment, error)
	// ProcessPayment...
}
