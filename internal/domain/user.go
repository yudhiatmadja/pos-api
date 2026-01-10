package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User Model (Mapped from DB or Custom)
type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Roles        []string  `json:"roles"`
	CreatedAt    time.Time `json:"created_at"`
}

// Auth Request/Response
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token,omitempty"` // simplified
	User         User      `json:"user"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"` // Creating initial user with role
}

// Interfaces

type UserRepository interface {
	CreateUser(ctx context.Context, username, passwordHash, role string) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	// Add other methods as needed
}

type AuthUsecase interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
}