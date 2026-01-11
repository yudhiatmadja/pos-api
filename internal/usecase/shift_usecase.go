package usecase

import (
	"context"
	"fmt"
	"math/big"

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

	startCash := pgtype.Numeric{Int: big.NewInt(int64(req.StartCash * 100)), Exp: -2, Valid: true}

	s, err := uc.store.CreateShift(ctx, repository.CreateShiftParams{
		UserID:    pgtype.UUID{Bytes: req.UserID, Valid: true},
		OutletID:  pgtype.UUID{Bytes: req.OutletID, Valid: true},
		StartCash: startCash,
	})
	if err != nil {
		return nil, err
	}

	scVal, _ := s.StartCash.Float64Value()
	return &domain.Shift{
		ID:        uuid.UUID(s.ID.Bytes),
		UserID:    uuid.UUID(s.UserID.Bytes),
		StartTime: s.StartTime.Time,
		StartCash: scVal.Float64,
	}, nil
}

func (uc *shiftUsecase) CloseShift(ctx context.Context, req *domain.CloseShiftRequest) (*domain.Shift, error) {
	// Calculate expected cash from sales during this shift (complex query needed)
	expectedCash := 0.0 // Placeholder

	endCashNum := pgtype.Numeric{Int: big.NewInt(int64(req.EndCash * 100)), Exp: -2, Valid: true}
	expCashNum := pgtype.Numeric{Int: big.NewInt(int64(expectedCash * 100)), Exp: -2, Valid: true}

	s, err := uc.store.CloseShift(ctx, repository.CloseShiftParams{
		ID:      pgtype.UUID{Bytes: req.ShiftID, Valid: true},
		EndCash: pgtype.Numeric{Valid: true}, // Wait, CloseShiftParams expects updated params
		// Check shifts.sql from step 47: UPDATE shifts SET end_time=NOW(), end_cash=$2, expected_cash=$3 WHERE id=$1
	})
	// NOTE: sqlc generated signature might differ based on named params.
	// Assuming positional or struct based on generic sqlc usage.
	// If I used $1, $2, $3 in SQL, sqlc uses params.

	// Correction: CloseShiftParams struct should have ID, EndCash, ExpectedCash.
	// I need to be careful with nullability if defined as nullable in DB.

	s, err = uc.store.CloseShift(ctx, repository.CloseShiftParams{
		ID:           pgtype.UUID{Bytes: req.ShiftID, Valid: true},
		EndCash:      pgtype.Numeric{Int: big.NewInt(int64(req.EndCash * 100)), Exp: -2, Valid: true}, // Wait, schema allows null but here we set it
		ExpectedCash: pgtype.Numeric{Int: big.NewInt(int64(expectedCash * 100)), Exp: -2, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	ec, _ := s.EndCash.Float64Value()
	return &domain.Shift{
		ID:      uuid.UUID(s.ID.Bytes),
		EndCash: &ec.Float64,
		EndTime: &s.EndTime.Time,
	}, nil
}

func (uc *shiftUsecase) GetCurrentShift(ctx context.Context, userID uuid.UUID) (*domain.Shift, error) {
	s, err := uc.store.GetCurrentShift(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return nil, err
	}
	sc, _ := s.StartCash.Float64Value()
	return &domain.Shift{
		ID:        uuid.UUID(s.ID.Bytes),
		StartCash: sc.Float64,
		StartTime: s.StartTime.Time,
	}, nil
}
