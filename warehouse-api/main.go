package main

import (
	"log"
	"net/http"
    "strings"
	"warehouse-api/config"
	"warehouse-api/handlers"
	"warehouse-api/middleware"
	"warehouse-api/repositories"
    "warehouse-api/services"
    
    _ "warehouse-api/docs" // Import swagger docs
    httpSwagger "github.com/swaggo/http-swagger"
)

// @title Warehouse Inventory API
// @version 1.0
// @description API for managing warehouse inventory, stock, and transactions.
// @host localhost:8080
// @BasePath /api

// @tag.name Auth
// @tag.description Autentikasi dan manajemen pengguna

// @tag.name Dashboard
// @tag.description Statistik dan ringkasan data gudang

// @tag.name Barang
// @tag.description Manajemen data barang inventaris

// @tag.name Stok
// @tag.description Manajemen dan monitoring stok barang

// @tag.name Pembelian
// @tag.description Transaksi pembelian dan stok masuk

// @tag.name Penjualan
// @tag.description Transaksi penjualan dan stok keluar

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// 1. Connect to Database
	config.ConnectDB()
    defer config.DB.Close()

	// 2. Initialize Repositories
	userRepo := repositories.NewUserRepository(config.DB)
	barangRepo := repositories.NewBarangRepository(config.DB)
	stokRepo := repositories.NewStokRepository(config.DB)
	pembelianRepo := repositories.NewPembelianRepository(config.DB)
    penjualanRepo := repositories.NewPenjualanRepository(config.DB)
    dashboardRepo := repositories.NewDashboardRepository(config.DB)

	// 3. Initialize Services
	userService := services.NewUserService(userRepo)
    penjualanService := services.NewPenjualanService(config.DB, penjualanRepo, stokRepo, barangRepo)
    pembelianService := services.NewPembelianService(config.DB, pembelianRepo, stokRepo, barangRepo)

	// 4. Initialize Handlers
	userHandler := handlers.NewUserHandler(userService)
	barangHandler := handlers.NewBarangHandler(barangRepo)
	stokHandler := handlers.NewStokHandler(stokRepo)
	pembelianHandler := handlers.NewPembelianHandler(pembelianService, pembelianRepo)
    penjualanHandler := handlers.NewPenjualanHandler(penjualanService, penjualanRepo)
    dashboardHandler := handlers.NewDashboardHandler(dashboardRepo)

	// 5. Setup Router
	mux := http.NewServeMux()

    // Swagger Route
    mux.Handle("/swagger/", httpSwagger.WrapHandler)

    // --- Routes Definition ---
    // Auth
	mux.HandleFunc("POST /api/login", userHandler.Login)
	mux.HandleFunc("POST /api/register", userHandler.Register)
    mux.HandleFunc("GET /api/users", userHandler.GetAll)

    // Barang
	mux.HandleFunc("GET /api/barang", barangHandler.GetAll)
    mux.HandleFunc("GET /api/barang/stok", barangHandler.GetAllWithStok)
	mux.HandleFunc("GET /api/barang/{id}", barangHandler.GetByID)
	mux.HandleFunc("POST /api/barang", barangHandler.Create)
	mux.HandleFunc("PUT /api/barang/{id}", barangHandler.Update)
	mux.HandleFunc("DELETE /api/barang/{id}", barangHandler.Delete)

    // Stok
	mux.HandleFunc("GET /api/stok", stokHandler.GetAll)
	mux.HandleFunc("GET /api/stok/{id}", stokHandler.GetByBarangID)
    mux.HandleFunc("GET /api/history-stok", stokHandler.GetHistory)
	mux.HandleFunc("GET /api/history-stok/{id}", stokHandler.GetHistory)

    // Pembelian
    mux.HandleFunc("POST /api/pembelian", pembelianHandler.Create)
    mux.HandleFunc("GET /api/pembelian", pembelianHandler.GetAll)
    mux.HandleFunc("GET /api/pembelian/{id}", pembelianHandler.GetByID)
    
    // Penjualan
    mux.HandleFunc("POST /api/penjualan", penjualanHandler.Create)
    mux.HandleFunc("GET /api/penjualan", penjualanHandler.GetAll)
    mux.HandleFunc("GET /api/penjualan/{id}", penjualanHandler.GetByID)

    // Dashboard
    mux.HandleFunc("GET /api/dashboard", dashboardHandler.GetStats)

    // --- Middleware Chains ---
    
    // 1. Auth Middleware Wrapper
    // Kita buat wrapper agar hanya route tertentu yang dicek auth-nya
    authHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Skip auth untuk login dan swagger
        if r.URL.Path == "/api/login" || strings.HasPrefix(r.URL.Path, "/swagger/") {
            mux.ServeHTTP(w, r)
            return
        }
        
        // Cek Auth untuk sisanya
        middleware.AuthMiddleware(mux).ServeHTTP(w, r)
    })

    // 2. CORS Handler
    corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        authHandler.ServeHTTP(w, r)
    })

    // 3. Logger & Rate Limit (Global)
    // Urutan: Logger -> RateLimit -> CORS -> Auth -> Mux
    finalHandler := middleware.Logger(middleware.RateLimitMiddleware(corsHandler))

	// 6. Start Server
	log.Println("Server starting on :8080")
    log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
