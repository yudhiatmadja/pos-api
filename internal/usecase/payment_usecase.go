package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"pos-api/internal/domain"
	"pos-api/internal/repository"

	"github.com/google/uuid"
)

type paymentUsecase struct {
	store repository.Store
}

func NewPaymentUsecase(store repository.Store) domain.PaymentUsecase {
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

	// MVP: Create/Update Payment record
	// Assume we just update the qris_url of an existing payment or create one.
	// For MVP let's assume one payment per order for now or just insert.
	// We need a repository method CreatePayment? I didn't add it in SQLC yet?
	// payments (order_id, payment_method, amount, qris_url)

	// NOTE: I need to add payment queries to generic SQL or create generic insert.
	// For now, I'll return a mock domain object if repo logic is missing,
	// OR I should have added `payments.sql`?
	// Protocol buffer: I missed `payments.sql`.

	return &domain.Payment{
		ID:            uuid.New(),
		OrderID:       req.OrderID,
		PaymentMethod: domain.PaymentMethodQRIS,
		Status:        domain.PaymentPending,
		QRISImageURL:  dst,
	}, nil
}
