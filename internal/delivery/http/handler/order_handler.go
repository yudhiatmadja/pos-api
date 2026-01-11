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

func NewOrderHandler(uc domain.OrderUsecase) *OrderHandler {
	return &OrderHandler{
		OrderUsecase: uc,
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

	// Get role from context (set by AuthMiddleware, standardizing on single "role" for logic check, or check all?)
	// For simplicity, we assume primary role or check if any of them is owner.
	// But AuthMiddleware sets "roles" []string.
	// Use helper to get roles.
	var userRole string
	if rolesVal, exists := c.Get("roles"); exists {
		if roles, ok := rolesVal.([]string); ok && len(roles) > 0 {
			userRole = roles[0] // Just taking first for now or we need to pass all?
			// Logic requires checking if ONE of them is StoreOwner.
			// Let's pass the first one, or better change usecase to accept slice?
			// Simpler: Check explicit permission here or pass strings.
			// Re-reading logic: "if userRole != StoreOwner".
			// If user has multiple roles, we should check if *any* is StoreOwner.
			for _, r := range roles {
				if r == string(domain.RoleStoreOwner) || r == string(domain.RoleSuperAdmin) {
					userRole = r
					break
				}
			}
			if userRole == "" {
				userRole = roles[0]
			}
		}
	}

	order, err := h.OrderUsecase.UpdateStatus(c.Request.Context(), orderID, domain.OrderStatus(req.Status), userID, userRole)
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
