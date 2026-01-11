package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TableSession struct {
	ID        uuid.UUID `json:"id"`
	TableID   uuid.UUID `json:"table_id"`
	StoreID   uuid.UUID `json:"store_id"`
	TableName string    `json:"table_name,omitempty"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IsActive  bool      `json:"is_active"`
}

type CreateSessionRequest struct {
	TableID uuid.UUID `json:"table_id" binding:"required"`
}

type SessionUsecase interface {
	CreateSession(ctx context.Context, tableID uuid.UUID) (*TableSession, error)
	ValidateSession(ctx context.Context, token string) (*TableSession, error)
}

type SessionRepository interface {
	// ... legacy if needed
}
