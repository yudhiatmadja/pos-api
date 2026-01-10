package v1

import (
	"net/http"

	"pos-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewPaymentHandler(router *gin.RouterGroup, uc domain.OrderUsecase) {
	handler := &PaymentHandler{
		orderUsecase: uc,
	}

	payments := router.Group("/payments")
	{
		payments.POST("/webhook", handler.handleWebhook)
	}
}

// handleWebhook simulates receiving a payment notification (e.g. Midtrans)
func (h *PaymentHandler) handleWebhook(ctx *gin.Context) {
	var req struct {
		OrderID           string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real system, we would verify signature first.
	// Then map external status to internal status.

	if req.TransactionStatus == "settlement" || req.TransactionStatus == "capture" {
		// Update order to PAID
		// We need method ReleaseOrder or UpdatePaymentStatus in Usecase
		// For now just logging or assuming success

		// This requires expanding OrderUsecase interface.
		// For MVP, we'll just acknowledge.
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
}
