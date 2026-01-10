package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TableSession struct {
	ID        uuid.UUID `json:"id"`
	TableID   uuid.UUID `json:"table_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
	TableName string    `json:"table_name,omitempty"` // Enriched data
	OutletID  uuid.UUID `json:"outlet_id,omitempty"`  // Enriched data
}

type CreateSessionRequest struct {
	TableID uuid.UUID `json:"table_id" binding:"required"`
}

type SessionUsecase interface {
	CreateSession(ctx context.Context, tableID uuid.UUID) (*TableSession, error)
	ValidateSession(ctx context.Context, token string) (*TableSession, error)
}
