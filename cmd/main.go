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

	// Adapting AuthUsecase to use 'store' queries (assuming it takes repository.Queries compatible interface)
	// If AuthUsecase expects *Queries, Store satisfies Querier interface (embedded) if initialized correctly.
	// My SQLStore embeds *Queries.
	// But `store` variable is interface `Store`.
	// Check auth_usecase.go if needed, but likely it uses generated interface.
	// For safety, I'll pass store.Queries (type assertion or method?)
	// SQLStore has *Queries embedded.
	// But store is interface. I need to expose Queries from Store?
	// Or just type assert.
	sqlStore := store.(*repository.SQLStore)

	authUsecase := usecase.NewAuthUsecase(sqlStore.Queries, connPool, tokenMaker, authConfig)
	sessionUsecase := usecase.NewSessionUsecase(sqlStore.Queries)

	// New Order Usecase with EventService
	orderUsecase := usecase.NewOrderUsecase(store, hub)

	// New Shift Usecase
	shiftUsecase := usecase.NewShiftUsecase(store)

	// 4. Setup Router
	router := gin.Default()
	router.Use(gin.Recovery())

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH") // Added PATCH
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Public Routes
	apiV1 := router.Group("/api/v1")
	v1.NewAuthHandler(apiV1, authUsecase)
	v1.NewSessionHandler(apiV1, sessionUsecase)

	// Auth Middleware for my handlers
	// middleware.AuthMiddleware expects util.TokenMaker
	authMiddleware := middleware.AuthMiddleware(tokenMaker)

	// Use my Handlers
	handler.NewOrderHandler(router, orderUsecase, authMiddleware)
	handler.NewShiftHandler(router, shiftUsecase, authMiddleware)
	// v1.NewPaymentHandler(apiV1, orderUsecase) // payment handled in OrderHandler potentially or need separate

	// WebSocket Route
	router.GET("/ws", func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	// 5. Start Server
	log.Printf("Starting server on %s", serverAddress)
	err = router.Run(serverAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
