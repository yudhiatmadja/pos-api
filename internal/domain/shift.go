package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	ID           uuid.UUID  `json:"id"`
	UserID       uuid.UUID  `json:"user_id"`
	OutletID     uuid.UUID  `json:"outlet_id"`
	StartTime    time.Time  `json:"start_time"`
	EndTime      *time.Time `json:"end_time,omitempty"`
	StartCash    float64    `json:"start_cash"`
	EndCash      *float64   `json:"end_cash,omitempty"`
	ExpectedCash *float64   `json:"expected_cash,omitempty"`
}

type OpenShiftRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	OutletID  uuid.UUID `json:"outlet_id"`
	StartCash float64   `json:"start_cash"`
}

type CloseShiftRequest struct {
	ShiftID uuid.UUID `json:"shift_id"`
	EndCash float64   `json:"end_cash"`
}

type ShiftUsecase interface {
	OpenShift(ctx context.Context, req *OpenShiftRequest) (*Shift, error)
	CloseShift(ctx context.Context, req *CloseShiftRequest) (*Shift, error)
	GetCurrentShift(ctx context.Context, userID uuid.UUID) (*Shift, error)
}
