package handler

import (
	"net/http"

	"pos-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderHandler struct {
	OrderUsecase domain.OrderUsecase
}

func NewOrderHandler(router *gin.Engine, uc domain.OrderUsecase, middleware gin.HandlerFunc) {
	handler := &OrderHandler{
		OrderUsecase: uc,
	}

	api := router.Group("/api/v1")
	{
		api.POST("/orders", handler.CreateOrder)
		api.GET("/orders/:id", handler.GetOrder)
		// Protected routes
		protected := api.Group("/", middleware)
		protected.PATCH("/orders/:id/status", handler.UpdateStatus)
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req domain.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.OrderUsecase.CreateOrder(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	order, err := h.OrderUsecase.UpdateStatus(c.Request.Context(), orderID, domain.OrderStatus(req.Status), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	idParam := c.Param("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.OrderUsecase.GetOrder(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
