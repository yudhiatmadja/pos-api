package usecase

import (
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"

	"pos-api/internal/domain"
	"pos-api/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
)

type paymentUsecase struct {
	store repository.Repository
}

func NewPaymentUsecase(store repository.Repository) domain.PaymentUsecase {
	return &paymentUsecase{store: store}
}

func (uc *paymentUsecase) UploadQRIS(ctx context.Context, req *domain.UploadQRISRequest) (*domain.Payment, error) {
	// 1. Validate Order exists?
	// 2. Save file
	// 3. Update/Create Payment record

	// Create upload dir
	uploadDir := "uploads/qris"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return nil, err
	}

	filename := fmt.Sprintf("%s-%s", req.OrderID, req.File.Filename)
	dst := filepath.Join(uploadDir, filename)

	src, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return nil, err
	}

	publicURL := fmt.Sprintf("/uploads/qris/%s", filename)

	// MVP: Create/Update Payment record
	// Assume we just update the qris_url of an existing payment or create one.
	// 4. Update Payment Record in DB
	// Assume request contains OrderID to link payment, or we treat ID as PaymentID??
	// The UploadQRISRequest in domain current has OrderID? Let's check.
	// If it doesn't, we might need to assume the ID passed is PaymentID or OrderID.
	// Let's assume for now we create a new payment record if it doesn't exist or update if it does.
	// But usually, QRIS is generated for an Order.
	// Let's create a payment record linked to OrderID.

	// Need to parse OrderID from request if available.
	// For MVP, let's assume we update the payment if we have a payment ID, or create one if we have order ID.
	// Domain struct: UploadQRISRequest { OrderID uuid.UUID, Image file ... }

	err = uc.store.ExecTx(ctx, func(q *repository.Queries) error {
		// Try to find existing payment for order
		existing, err := q.GetPaymentByOrder(ctx, pgtype.UUID{Bytes: req.OrderID, Valid: true})
		if err != nil {
			// Create new payment
			amount := pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true} // Unknown amount yet?
			_, err = q.CreatePayment(ctx, repository.CreatePaymentParams{
				OrderID:       pgtype.UUID{Bytes: req.OrderID, Valid: true},
				PaymentMethod: "QRIS",
				Amount:        amount, // Placeholder
				Status:        "PENDING",
				QrisUrl:       pgtype.Text{String: publicURL, Valid: true},
			})
			return err
		} else {
			// Update existing
			err = q.UpdatePaymentQRIS(ctx, repository.UpdatePaymentQRISParams{
				ID:      existing.ID,
				QrisUrl: pgtype.Text{String: publicURL, Valid: true},
			})
			return err
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to save payment record: %w", err)
	}

	return &domain.Payment{
		OrderID:      req.OrderID,
		QRISImageURL: publicURL,
		Status:       domain.PaymentPending,
	}, nil
}
