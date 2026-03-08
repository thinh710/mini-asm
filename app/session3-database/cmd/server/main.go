package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"mini-asm/internal/database"
	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/postgres"

	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
	log.Println("🚀 Starting Mini ASM Server (Session 3 - Database)...")

	// ============================================
	// CONFIGURATION - Load from environment
	// ============================================

	// Database configuration
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "secops")
	dbPassword := getEnv("DB_PASSWORD", "secops123")
	dbName := getEnv("DB_NAME", "mini_asm")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode,
	)

	log.Printf("📊 Connecting to database: %s@%s:%s/%s", dbUser, dbHost, dbPort, dbName)

	// ============================================
	// DATABASE CONNECTION
	// ============================================

	// Open database connection
	db, err := database.ConnectWithRetry(connStr, 5)
	if err != nil {
		log.Fatalf("❌ Could not connect to database after retries: %v", err)
	}
	defer db.Close()

	// Optional: Configure connection pool
	db.SetMaxOpenConns(25)               // Maximum open connections
	db.SetMaxIdleConns(5)                // Maximum idle connections
	db.SetConnMaxLifetime(5 * 60 * 1000) // Connection lifetime (5 minutes)

	// ============================================
	// DEPENDENCY INJECTION - Wire up all layers
	// ============================================

	// 1. Initialize Storage Layer (Infrastructure)
	//    🎯 KEY CHANGE: PostgresStorage instead of MemoryStorage!
	//    Compare with Session 2:
	//    Session 2: store := memory.NewMemoryStorage()
	//    Session 3: store := postgres.NewPostgresStorage(db)
	store := postgres.NewPostgresStorage(db)
	log.Println("✅ Storage initialized: PostgreSQL")

	// 2. Initialize Service Layer (Use Case / Business Logic)
	//    ✨ NO CHANGES! Service doesn't care about storage implementation
	assetService := service.NewAssetService(store)
	log.Println("✅ Service initialized: AssetService")

	// 3. Initialize Handler Layer (Presentation / HTTP)
	//    ✨ NO CHANGES! Handler doesn't care about storage implementation
	assetHandler := handler.NewAssetHandler(assetService)
	healthHandler := handler.NewHealthHandler(db)
	log.Println("✅ Handlers initialized")

	// ============================================
	// ROUTING - Register HTTP endpoints
	// ============================================

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", healthHandler.Check)
	// BÀi 1
	mux.HandleFunc("GET /assets/stats", assetHandler.GetStats)
	mux.HandleFunc("GET /assets/count", assetHandler.CountAssets)
	// bài 2
	mux.HandleFunc("POST /assets/batch", assetHandler.BatchCreate)
	//bài 3
	mux.HandleFunc("DELETE /assets/batch", assetHandler.BatchDelete)
	//bài 7
	mux.HandleFunc("GET /assets/search", assetHandler.SearchAssets)

	// Asset CRUD operations
	mux.HandleFunc("POST /assets", assetHandler.CreateAsset)        // Create
	mux.HandleFunc("GET /assets", assetHandler.ListAssets)          // Read (list with filters)
	mux.HandleFunc("GET /assets/{id}", assetHandler.GetAsset)       // Read (single)
	mux.HandleFunc("PUT /assets/{id}", assetHandler.UpdateAsset)    // Update
	mux.HandleFunc("DELETE /assets/{id}", assetHandler.DeleteAsset) // Delete

	log.Println("✅ Routes registered:")
	log.Println("   GET    /health")
	log.Println("   POST   /assets")
	log.Println("   GET    /assets")
	log.Println("   GET    /assets/{id}")
	log.Println("   PUT    /assets/{id}")
	log.Println("   DELETE /assets/{id}")

	// ============================================
	// START SERVER
	// ============================================

	port := getEnv("SERVER_PORT", "8080")
	addr := ":" + port

	log.Printf("🌐 Server listening on http://localhost%s\n", addr)
	log.Println("📖 API Documentation: see docs/api.yml")
	log.Println("🗄️  Database: PostgreSQL (persistent storage)")
	log.Println("Press Ctrl+C to stop")
	log.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

// getEnv retrieves an environment variable with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
