package usecase

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"pos-api/internal/domain"
	"pos-api/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type shiftUsecase struct {
	store repository.Store
}

func NewShiftUsecase(store repository.Store) domain.ShiftUsecase {
	return &shiftUsecase{store: store}
}

func (uc *shiftUsecase) OpenShift(ctx context.Context, req *domain.OpenShiftRequest) (*domain.Shift, error) {
	// Check if user already has open shift
	existing, err := uc.store.GetCurrentShift(ctx, pgtype.UUID{Bytes: req.UserID, Valid: true})
	if err == nil && existing.ID.Valid {
		return nil, fmt.Errorf("user already has an active shift")
	}

	startCash := pgtype.Numeric{Int: big.NewInt(int64(req.OpeningCash * 100)), Exp: -2, Valid: true}

	s, err := uc.store.CreateShift(ctx, repository.CreateShiftParams{
		UserID:      pgtype.UUID{Bytes: req.UserID, Valid: true},
		StoreID:     pgtype.UUID{Bytes: req.StoreID, Valid: true},
		OpeningCash: startCash,
	})
	if err != nil {
		return nil, err
	}

	scVal, _ := s.OpeningCash.Float64Value()
	return &domain.Shift{
		ID:          uuid.UUID(s.ID.Bytes),
		UserID:      uuid.UUID(s.UserID.Bytes),
		StoreID:     uuid.UUID(s.StoreID.Bytes),
		OpenedAt:    s.OpenedAt.Time,
		OpeningCash: scVal.Float64,
	}, nil
}

func (uc *shiftUsecase) CloseShift(ctx context.Context, req *domain.CloseShiftRequest) (*domain.Shift, error) {
	// Calculate expected cash from sales during this shift (complex query needed)
	expectedCash := 0.0 // Placeholder

	s, err := uc.store.CloseShift(ctx, repository.CloseShiftParams{
		ID:           pgtype.UUID{Bytes: req.ShiftID, Valid: true},
		ClosingCash:  pgtype.Numeric{Int: big.NewInt(int64(req.ClosingCash * 100)), Exp: -2, Valid: true},
		ExpectedCash: pgtype.Numeric{Int: big.NewInt(int64(expectedCash * 100)), Exp: -2, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	ec, _ := s.ClosingCash.Float64Value()
	exp, _ := s.ExpectedCash.Float64Value()

	// Handle nullable ClosedAt
	var closedAt *time.Time
	if s.ClosedAt.Valid {
		t := s.ClosedAt.Time
		closedAt = &t
	}

	return &domain.Shift{
		ID:           uuid.UUID(s.ID.Bytes),
		UserID:       uuid.UUID(s.UserID.Bytes),
		StoreID:      uuid.UUID(s.StoreID.Bytes),
		OpenedAt:     s.OpenedAt.Time, // We might need to fetch this or just return updated fields
		ClosingCash:  &ec.Float64,
		ExpectedCash: &exp.Float64,
		ClosedAt:     closedAt,
	}, nil
}

func (uc *shiftUsecase) GetCurrentShift(ctx context.Context, userID uuid.UUID) (*domain.Shift, error) {
	s, err := uc.store.GetCurrentShift(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	sc, _ := s.OpeningCash.Float64Value()
	return &domain.Shift{
		ID:          uuid.UUID(s.ID.Bytes),
		UserID:      uuid.UUID(s.UserID.Bytes),
		StoreID:     uuid.UUID(s.StoreID.Bytes),
		OpenedAt:    s.OpenedAt.Time,
		OpeningCash: sc.Float64,
	}, nil
}
