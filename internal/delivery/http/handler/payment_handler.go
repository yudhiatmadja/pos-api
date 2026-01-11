package handler

import (
	"net/http"

	"pos-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	PaymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(router *gin.Engine, uc domain.PaymentUsecase, middleware gin.HandlerFunc) {
	handler := &PaymentHandler{
		PaymentUsecase: uc,
	}

	// Group: /api/v1/payments
	api := router.Group("/api/v1/payments")
	// Protected
	protected := api.Group("/", middleware)
	{
		protected.POST("/qris/upload", handler.UploadQRIS)
	}
}

func (h *PaymentHandler) UploadQRIS(c *gin.Context) {
	// Multipart form
	// order_id: uuid
	// image: file

	orderIDStr := c.PostForm("order_id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing order_id"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	req := &domain.UploadQRISRequest{
		OrderID: orderID,
		File:    file,
	}

	payment, err := h.PaymentUsecase.UploadQRIS(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, payment)
}
