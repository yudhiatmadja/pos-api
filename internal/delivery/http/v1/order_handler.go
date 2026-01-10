package v1

import (
	"net/http"

	"pos-api/internal/delivery/http/middleware"
	"pos-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderHandler(router *gin.RouterGroup, uc domain.OrderUsecase) {
	handler := &OrderHandler{
		orderUsecase: uc,
	}

	orders := router.Group("/orders")
	{
		// Order creation can be public (QR) or protected (Cashier)
		// For simplicity, we assume the token (either User JWT or Session Token)
		// is passed in headers and validated by logic or middleware.
		// Here we allow it to be open but Usecase should validate ownership?
		// Better: specific endpoints or unified with logic.
		// For now: Open endpoint for creating order, but we should probably
		// check if session exists if passed.
		orders.POST("", handler.createOrder)
	}
}

func (h *OrderHandler) createOrder(ctx *gin.Context) {
	var req domain.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Logic to inject User ID if authenticated via JWT (Cashier)
	payload, exists := ctx.Get(middleware.AuthorizationPayloadKey)
	if exists {
		// handle cashier logic if needed
		_ = payload
	}

	order, err := h.orderUsecase.CreateOrder(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, order)
}
