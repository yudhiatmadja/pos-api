package domain

import "context"

type EventService interface {
	PublishEvent(ctx context.Context, eventType string, payload interface{}) error
}
