package v1

import (
	"net/http"

	"pos-api/internal/domain"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionUsecase domain.SessionUsecase
}

func NewSessionHandler(router *gin.RouterGroup, uc domain.SessionUsecase) {
	handler := &SessionHandler{
		sessionUsecase: uc,
	}

	sessions := router.Group("/table-sessions")
	{
		sessions.POST("", handler.createSession)
		// Validation is usually done via Middleware or explicit check
		sessions.POST("/validate", handler.validateSession)
	}
}

func (h *SessionHandler) createSession(ctx *gin.Context) {
	var req domain.CreateSessionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionUsecase.CreateSession(ctx, req.TableID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) validateSession(ctx *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, err := h.sessionUsecase.ValidateSession(ctx, req.Token)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, session)
}
