package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"pos-api/internal/delivery/http/middleware"
	v1 "pos-api/internal/delivery/http/v1"
	"pos-api/internal/delivery/ws"
	"pos-api/internal/repository"
	"pos-api/internal/usecase"
	"pos-api/internal/util"
)

func main() {
	// 1. Config (Hardcoded for now, ideal to load from env)
	dbSource := "postgresql://root:secret@localhost:5432/pos_db?sslmode=disable"
	serverAddress := "0.0.0.0:8080"
	tokenSymmetricKey := "12345678901234567890123456789012" // 32 chars
	accessTokenDuration := 24 * time.Hour

	// 2. Setup Database
	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer connPool.Close()

	// Redis Setup
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	hub := ws.NewHub(redisClient)
	go hub.Run()

	store := repository.New(connPool)

	// 3. Setup Dependencies
	tokenMaker, err := util.NewJWTMaker(tokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	authConfig := usecase.AuthConfig{
		AccessTokenDuration: accessTokenDuration,
	}
	authUsecase := usecase.NewAuthUsecase(store, connPool, tokenMaker, authConfig)
	sessionUsecase := usecase.NewSessionUsecase(store)
	orderUsecase := usecase.NewOrderUsecase(store, connPool)

	// 4. Setup Router
	router := gin.Default()
	router.Use(gin.Recovery())

	// CORS Middleware (Simple version)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	apiV1 := router.Group("/api/v1")

	// Public Routes
	v1.NewAuthHandler(apiV1, authUsecase)
	v1.NewSessionHandler(apiV1, sessionUsecase)
	v1.NewOrderHandler(apiV1, orderUsecase)
	v1.NewPaymentHandler(apiV1, orderUsecase)

	// WebSocket Route
	apiV1.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// Protected Routes
	protected := apiV1.Group("/")
	protected.Use(middleware.AuthMiddleware(tokenMaker))
	{
		protected.GET("/ping", func(ctx *gin.Context) {
			payload, _ := ctx.Get(middleware.AuthorizationPayloadKey)
			ctx.JSON(200, gin.H{
				"message": "pong",
				"user":    payload,
			})
		})
	}

	// 5. Start Server
	log.Printf("Starting server on %s", serverAddress)
	err = router.Run(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
