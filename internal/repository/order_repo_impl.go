package repository

import (
	"context"

	"pos-api/internal/domain"
)

type OrderRepositoryImpl struct {
	q *Queries
}

func NewOrderRepository(q *Queries) *OrderRepositoryImpl {
	return &OrderRepositoryImpl{q: q}
}

// Map methods...

func (r *OrderRepositoryImpl) CreateOrder(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	// Transaction support is usually needed here for Order + OrderItems.
	// For now implementing single table insert wrapper, but UseCase should handle TX.
	// We will expose a method to run in TX in the Store or similar.
	return nil, nil // Placeholder, implemented in Usecase likely via Store pattern
}

// Instead of individual repos, usually with SQLC and TX, we use a Store struct.
// I will create a Store struct that embeds Queries and handles transactions.
