package main

import (
	"log"
	"net/http"

	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/memory"
)

func main() {
	log.Println("🚀 Starting Mini ASM Server...")

	// ============================================
	// DEPENDENCY INJECTION - Wire up all layers
	// ============================================

	// 1. Initialize Storage Layer (Infrastructure)
	//    Using in-memory storage for now
	//    Buổi 3 sẽ swap sang PostgreSQL - chỉ thay đổi dòng này!
	store := memory.NewMemoryStorage()
	log.Println("✅ Storage initialized: In-Memory")

	// 2. Initialize Service Layer (Use Case / Business Logic)
	//    Inject storage dependency
	assetService := service.NewAssetService(store)
	log.Println("✅ Service initialized: AssetService")

	// 3. Initialize Handler Layer (Presentation / HTTP)
	//    Inject service dependency
	assetHandler := handler.NewAssetHandler(assetService)
	healthHandler := handler.NewHealthHandler()
	log.Println("✅ Handlers initialized")

	// ============================================
	// ROUTING - Register HTTP endpoints
	// ============================================

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", healthHandler.Check)

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

	addr := ":8080"
	log.Printf("🌐 Server listening on http://localhost%s\n", addr)
	log.Println("📖 API Documentation: see docs/api.yml")
	log.Println("Press Ctrl+C to stop")
	log.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

/*

=== HOMEWORK HINTS ===

Students có thể:
1. Add logging middleware
2. Add request ID tracking
3. Add CORS headers
4. Add rate limiting
5. Add metrics endpoint

Tất cả là middleware, không touch business logic!
*/
