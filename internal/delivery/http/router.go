package http

import (
	"github.com/gin-gonic/gin"
	"pos-api/internal/delivery/http/handler"
	"pos-api/internal/delivery/http/middleware"
)

func NewRouter(r *gin.Engine, authHandler *handler.AuthHandler) {
	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Contoh Protected Route
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware()) // Pasang middleware di sini
		{
			admin.GET("/dashboard", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Welcome Admin"})
			})
		}
	}
}