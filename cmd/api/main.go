package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	
	"pos-api/internal/delivery/http"
	"pos-api/internal/delivery/http/handler"
	"pos-api/internal/repository"
	"pos-api/internal/usecase"
	"pos-api/pkg/database"
)

func main() {
	// 1. Load Config
	dbSource := os.Getenv("DB_SOURCE")
	if dbSource == "" {
		// Connect ke PostgreSQL di Docker (dari aplikasi Go yang jalan di local)
		// Gunakan localhost karena port 5432 sudah di-expose ke host
		dbSource = "postgresql://postgres:password@localhost:5432/mypos_db?sslmode=disable"
		log.Println("⚠️  Using default DB connection")
	}

	// 2. Connect Database
	log.Println("🔌 Connecting to database...")
	dbPool, err := database.NewPostgresConnection(dbSource)
	if err != nil {
		log.Fatalf("❌ Cannot connect to db: %v\n\nPastikan PostgreSQL container running:\n  docker-compose up -d postgres\n\n", err)
	}
	defer dbPool.Close()
	
	log.Println("✅ Database connected successfully")

	// 3. Setup Layers (Dependency Injection)
	repo := repository.New(dbPool)
	authUC := usecase.NewAuthUsecase(repo)
	authHandler := handler.NewAuthHandler(authUC)

	// 4. Setup Router
	r := gin.Default()
	http.NewRouter(r, authHandler)

	// 5. Run Server
	log.Println("🚀 Server running on http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}