package usecase

import (
	"context"
	"fmt"
	"time"

	"pos-api/internal/domain"
	"pos-api/internal/repository"
	"pos-api/internal/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type authUsecase struct {
	repo       *repository.Queries // SQLC generated queries
	pool       *pgxpool.Pool       // For transactions if needed
	tokenMaker util.TokenMaker
	config     AuthConfig
}

type AuthConfig struct {
	AccessTokenDuration time.Duration
}

func NewAuthUsecase(repo *repository.Queries, pool *pgxpool.Pool, tokenMaker util.TokenMaker, config AuthConfig) domain.AuthUsecase {
	return &authUsecase{
		repo:       repo,
		pool:       pool,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create User
	// Note: We need to adapt CreateUser params based on generated code.
	// Assuming CreateUser takes (ctx, arg)

	// Transaction could be better here to ensure user + role assignment
	// But simply:
	user, err := uc.repo.CreateUser(ctx, repository.CreateUserParams{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         req.Role, // Legacy column
	})
	if err != nil {
		return nil, err
	}

	// Create proper Role entry
	err = uc.repo.AssignRoleToUser(ctx, repository.AssignRoleToUserParams{
		UserID:   user.ID,
		RoleCode: req.Role,
	})
	if err != nil {
		// Cleanup user if this fails? For now just return error
		return nil, fmt.Errorf("failed to assign role: %w", err)
	}

	return &domain.User{
		ID:        uuid.UUID(user.ID.Bytes),
		Username:  user.Username,
		Roles:     []string{user.Role},
		CreatedAt: user.CreatedAt.Time,
	}, nil
}

func (uc *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := uc.repo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	err = util.CheckPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Fetch roles
	// In the legacy schema we have 'role' column, but we also have user_roles table now.
	// Let's get roles from user_roles
	rolesRaw, err := uc.repo.GetUserRoles(ctx, user.ID)
	var roles []string
	if err != nil || len(rolesRaw) == 0 {
		// Fallback to single role column if no user_roles found (legacy support)
		roles = []string{user.Role}
	} else {
		for _, r := range rolesRaw {
			roles = append(roles, r.Code)
		}
	}

	accessToken, err := uc.tokenMaker.CreateToken(
		uuid.UUID(user.ID.Bytes),
		user.Username,
		roles,
		uc.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		AccessToken: accessToken,
		User: domain.User{
			ID:        uuid.UUID(user.ID.Bytes),
			Username:  user.Username,
			Roles:     roles,
			CreatedAt: user.CreatedAt.Time,
		},
	}, nil
}
