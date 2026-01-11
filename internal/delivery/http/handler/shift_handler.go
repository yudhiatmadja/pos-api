package handler

import (
	"net/http"

	"pos-api/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShiftHandler struct {
	ShiftUsecase domain.ShiftUsecase
}

func NewShiftHandler(router *gin.Engine, uc domain.ShiftUsecase, middleware gin.HandlerFunc) {
	handler := &ShiftHandler{
		ShiftUsecase: uc,
	}

	protected := router.Group("/api/v1/shifts", middleware)
	{
		protected.POST("/open", handler.OpenShift)
		protected.POST("/close", handler.CloseShift)
		protected.GET("/current", handler.GetCurrentShift)
	}
}

func (h *ShiftHandler) OpenShift(c *gin.Context) {
	var req domain.OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Override UserID from token if needed for security,
	// or validate req.UserID matches token.
	authUserIDStr := c.GetString("user_id")
	if authUserIDStr != "" && authUserIDStr != req.UserID.String() {
		// Optional: enforce same user
	}

	shift, err := h.ShiftUsecase.OpenShift(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, shift)
}

func (h *ShiftHandler) CloseShift(c *gin.Context) {
	var req domain.CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shift, err := h.ShiftUsecase.CloseShift(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, shift)
}

func (h *ShiftHandler) GetCurrentShift(c *gin.Context) {
	// user_id from query or token
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	shift, err := h.ShiftUsecase.GetCurrentShift(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active shift found"})
		return
	}

	c.JSON(http.StatusOK, shift)
}
