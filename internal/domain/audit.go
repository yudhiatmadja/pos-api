package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID              `json:"id"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Action    string                 `json:"action"`
	Entity    string                 `json:"entity,omitempty"`
	EntityID  *uuid.UUID             `json:"entity_id,omitempty"`
	Before    map[string]interface{} `json:"before,omitempty"`
	After     map[string]interface{} `json:"after,omitempty"`
	IPAddress string                 `json:"ip_address,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
}

type AuditUsecase interface {
	Log(ctx context.Context, log AuditLog) error
}
