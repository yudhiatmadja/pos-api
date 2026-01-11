package repository

import (
	"context"

	"pos-api/internal/domain"
)

// Ensure implementation
var _ domain.UserRepository = (*UserRepositoryImpl)(nil)

type UserRepositoryImpl struct {
	q *Queries
}

func NewUserRepository(q *Queries) *UserRepositoryImpl {
	return &UserRepositoryImpl{q: q}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, username, passwordHash, role string) (domain.User, error) {
	row, err := r.q.CreateUser(ctx, CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
		Role:         role,
	})
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:           row.ID.Bytes,
		Username:     row.Username,
		PasswordHash: passwordHash,
		Roles:        []string{row.Role}, // Simplified one role for now
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}

func (r *UserRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (domain.User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}
	// Note: sqlc query GetUserByUsername currently defined to return 'User' struct from models.go?
	// or specific row. I need to check users.sql.go again to see what it returns.
	// Assuming it returns a struct with fields.
	// Wait, earlier view_file showed it returns `User` which is `models.go` struct.

	return domain.User{
		ID:           row.ID.Bytes,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		Roles:        []string{row.Role},
		CreatedAt:    row.CreatedAt.Time,
	}, nil
}
