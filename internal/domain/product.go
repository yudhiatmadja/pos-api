package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID  `json:"id"`
	StoreID     uuid.UUID  `json:"store_id"`
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	SKU         string     `json:"sku,omitempty"`
	Price       float64    `json:"price"`
	Stock       int32      `json:"stock"`
	ImageURL    string     `json:"image_url,omitempty"`
	IsAvailable bool       `json:"is_available"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type CreateProductRequest struct {
	StoreID     uuid.UUID  `json:"store_id" binding:"required"`
	CategoryID  *uuid.UUID `json:"category_id"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	SKU         string     `json:"sku"`
	Price       float64    `json:"price" binding:"required,gt=0"`
	Stock       int32      `json:"stock" binding:"gte=0"`
	ImageURL    string     `json:"image_url"`
}

type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*Product, error)
	ListProducts(ctx context.Context, storeID uuid.UUID, page, limit int32) ([]Product, error)
}
