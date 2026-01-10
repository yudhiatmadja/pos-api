package usecase

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"pos-api/internal/domain"
	"pos-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type sessionUsecase struct {
	repo *repository.Queries
}

func NewSessionUsecase(repo *repository.Queries) domain.SessionUsecase {
	return &sessionUsecase{repo: repo}
}

func (uc *sessionUsecase) CreateSession(ctx context.Context, tableID uuid.UUID) (*domain.TableSession, error) {
	// Generate random token
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	token := hex.EncodeToString(bytes)

	// Set expiry (e.g., 2 hours)
	expiresAt := time.Now().Add(2 * time.Hour)

	// Create session in DB
	session, err := uc.repo.CreateSession(ctx, repository.CreateSessionParams{
		TableID:   tableID,
		Token:     token,
		ExpiresAt: pgtype.Timestamptz{Time: expiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return &domain.TableSession{
		ID:        session.ID,
		TableID:   session.TableID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt.Time,
		IsActive:  session.IsActive.Bool,
	}, nil
}

func (uc *sessionUsecase) ValidateSession(ctx context.Context, token string) (*domain.TableSession, error) {
	// Custom query needed for Join, or separate calls.
	// We added GetSessionByToken in sessions.sql which joins tables.
	// Expected generated signature: GetSessionByToken(ctx, token) (GetSessionByTokenRow, error)

	row, err := uc.repo.GetSessionByToken(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("invalid or expired session")
	}

	return &domain.TableSession{
		ID:        row.ID,
		TableID:   row.TableID,
		Token:     row.Token,
		ExpiresAt: row.ExpiresAt.Time,
		IsActive:  row.IsActive.Bool,
		TableName: row.TableName,
		OutletID:  row.OutletID,
	}, nil
}
