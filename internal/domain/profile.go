package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "super_admin"
	RoleStoreOwner UserRole = "store_owner"
	RoleStaff      UserRole = "staff"
	RoleSupplier   UserRole = "supplier"
)

type Profile struct {
	ID        uuid.UUID  `json:"id"`
	Email     string     `json:"email"`
	FullName  string     `json:"full_name"`
	Role      UserRole   `json:"role"`
	StoreID   *uuid.UUID `json:"store_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type RegisterRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=6"`
	FullName string   `json:"full_name" binding:"required"`
	Role     UserRole `json:"role" binding:"required"`
	StoreID  string   `json:"store_id"` // Optional
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	Profile     *Profile `json:"profile"`
}

type AuthUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*Profile, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}
