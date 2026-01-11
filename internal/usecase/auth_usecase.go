package usecase

import (
	"context"
	"fmt"
	"time"

	"pos-api/internal/domain"
	"pos-api/internal/repository"
	"pos-api/internal/util"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type authUsecase struct {
	store      repository.Repository
	tokenMaker util.TokenMaker
	config     AuthConfig
}

type AuthConfig struct {
	AccessTokenDuration time.Duration
}

func NewAuthUsecase(store repository.Repository, tokenMaker util.TokenMaker, config AuthConfig) domain.AuthUsecase {
	return &authUsecase{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
}

func (uc *authUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.Profile, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	var profile *domain.Profile

	err = uc.store.ExecTx(ctx, func(q *repository.Queries) error {
		// 1. Create Auth User (Simulated)
		resAuth, err := q.CreateAuthUser(ctx, repository.CreateAuthUserParams{
			Email:             req.Email,
			EncryptedPassword: hashedPassword,
		})
		if err != nil {
			return fmt.Errorf("failed to create auth user: %w", err)
		}

		// 2. Create Profile
		// Profile ID matches Auth User ID
		var storeID pgtype.UUID
		if req.StoreID != "" {
			parsedID, err := uuid.Parse(req.StoreID)
			if err == nil {
				storeID = pgtype.UUID{Bytes: parsedID, Valid: true}
			}
		}

		resProfile, err := q.CreateProfile(ctx, repository.CreateProfileParams{
			ID:       resAuth.ID, // Link ID
			Email:    pgtype.Text{String: req.Email, Valid: true},
			FullName: pgtype.Text{String: req.FullName, Valid: true},
			Role:     string(req.Role),
			StoreID:  storeID,
		})
		if err != nil {
			return fmt.Errorf("failed to create profile: %w", err)
		}

		// Map to domain
		var sID *uuid.UUID
		if resProfile.StoreID.Valid {
			uid := uuid.UUID(resProfile.StoreID.Bytes)
			sID = &uid
		}

		profile = &domain.Profile{
			ID:        uuid.UUID(resProfile.ID.Bytes),
			Email:     resProfile.Email.String,
			FullName:  resProfile.FullName.String,
			Role:      domain.UserRole(resProfile.Role),
			StoreID:   sID,
			CreatedAt: resProfile.CreatedAt.Time,
			// UpdatedAt: resProfile.UpdatedAt.Time, // Not available in DB yet
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return profile, nil
}

func (uc *authUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	// 1. Get Auth User
	authUser, err := uc.store.GetAuthUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials") // User not found
	}

	// 2. Check Password
	err = util.CheckPassword(req.Password, authUser.EncryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// 3. Get Profile
	profileDB, err := uc.store.GetProfile(ctx, authUser.ID)
	if err != nil {
		return nil, fmt.Errorf("profile not found")
	}

	// 4. Create Token
	// Role is in profile
	role := profileDB.Role

	accessToken, err := uc.tokenMaker.CreateToken(
		uuid.UUID(authUser.ID.Bytes),
		authUser.Email,
		[]string{role}, // Pass as slice
		uc.config.AccessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	// Map to domain
	var sID *uuid.UUID
	if profileDB.StoreID.Valid {
		uid := uuid.UUID(profileDB.StoreID.Bytes)
		sID = &uid
	}

	return &domain.LoginResponse{
		AccessToken: accessToken,
		Profile: &domain.Profile{
			ID:        uuid.UUID(profileDB.ID.Bytes),
			Email:     profileDB.Email.String,
			FullName:  profileDB.FullName.String,
			Role:      domain.UserRole(profileDB.Role),
			StoreID:   sID,
			CreatedAt: profileDB.CreatedAt.Time,
			// UpdatedAt: profileDB.UpdatedAt.Time,
		},
	}, nil
}
