package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"pos-api/internal/delivery/http/handler"
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

	// Use NewStore (Transaction support)
	store := repository.NewStore(connPool)

	// 3. Setup Dependencies
	tokenMaker, err := util.NewJWTMaker(tokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	authConfig := usecase.AuthConfig{
		AccessTokenDuration: accessTokenDuration,
	}

	// Usecases
	authUsecase := usecase.NewAuthUsecase(store, tokenMaker, authConfig)

	// SessionUsecase likely uses *repository.Queries (need to verify its impl in existing codebase or assume)
	// If it hasn't been updated to Store, we extract Queries.
	sqlStore := store.(*repository.SQLStore)
	sessionUsecase := usecase.NewSessionUsecase(sqlStore.Queries)

	orderUsecase := usecase.NewOrderUsecase(store, hub)
	shiftUsecase := usecase.NewShiftUsecase(store)
	paymentUsecase := usecase.NewPaymentUsecase(store)

	// 4. Setup Router
	router := gin.Default()
	router.Use(gin.Recovery())

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// 5. Handlers & Routes
	apiV1 := router.Group("/api/v1")
	v1.NewAuthHandler(apiV1, authUsecase)
	v1.NewSessionHandler(apiV1, sessionUsecase)

	// Auth Middleware
	authMiddleware := middleware.AuthMiddleware(tokenMaker)

	// Protected Routes handled by specific Handlers
	handler.NewOrderHandler(router, orderUsecase, authMiddleware)
	handler.NewShiftHandler(router, shiftUsecase, authMiddleware)
	handler.NewPaymentHandler(router, paymentUsecase, authMiddleware)

	// WebSocket Route
	apiV1.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// 6. Start Server
	log.Printf("Starting server on %s", serverAddress)
	err = router.Run(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
