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

	// Validate user_id matches token
	authUserIDStr := c.GetString("user_id")
	if authUserIDStr != "" {
		// Just parse it to ensure validity,
		// Use it if req.UserID is empty, or enforce match.
		// MVP: Set it from token.
		uid, err := uuid.Parse(authUserIDStr)
		if err == nil {
			req.UserID = uid
		}
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
	userIDStr := c.GetString("user_id")
	userID, _ := uuid.Parse(userIDStr)

	shift, err := h.ShiftUsecase.GetCurrentShift(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active shift found"})
		return
	}

	c.JSON(http.StatusOK, shift)
}
