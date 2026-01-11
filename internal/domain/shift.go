package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	StoreID      uuid.UUID  `json:"store_id"`
	OpenedAt     time.Time  `json:"opened_at"`
	ClosedAt     *time.Time `json:"closed_at,omitempty"`
	OpeningCash  float64    `json:"opening_cash"`
	ClosingCash  *float64   `json:"closing_cash,omitempty"`
	ExpectedCash *float64   `json:"expected_cash,omitempty"`
}

type OpenShiftRequest struct {
	UserID      uuid.UUID `json:"user_id"` // Usually from token, or implicit
	StoreID     uuid.UUID `json:"store_id" binding:"required"`
	OpeningCash float64   `json:"opening_cash" binding:"required,gte=0"`
}

type CloseShiftRequest struct {
	ShiftID     uuid.UUID `json:"shift_id" binding:"required"`
	ClosingCash float64   `json:"closing_cash" binding:"required,gte=0"`
}

type ShiftUsecase interface {
	OpenShift(ctx context.Context, req *OpenShiftRequest) (*Shift, error)
	CloseShift(ctx context.Context, req *CloseShiftRequest) (*Shift, error)
	GetCurrentShift(ctx context.Context, userID uuid.UUID) (*Shift, error)
}
