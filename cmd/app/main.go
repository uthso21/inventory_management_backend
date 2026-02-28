package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/database"
	"github.com/uthso21/inventory_management_backend/internal/middleware"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {

	// Load .env file (ignore error if not present — env vars may be set externally)
	_ = godotenv.Load()

	// =========================
	// Connect to Database FIRST
	// =========================
	database.Connect()

	// =========================
	// Wire Dependencies
	// =========================

	// User
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := httpHandler.NewUserHandler(userService)
	authHandler := httpHandler.NewAuthHandler(userService)

	// Warehouse
	warehouseRepo := repository.NewWarehouseRepository()
	warehouseService := service.NewWarehouseService(warehouseRepo)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)

	// Purchase
	purchaseRepo := repository.NewPurchaseRepository()
	purchaseService := service.NewPurchaseService(purchaseRepo)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)

	// ML Agent
	mlService := service.NewMLAgentServiceWithDefaults()
	mlHandler := httpHandler.NewMLAgentHandler(mlService)

	// =========================
	// Role shorthand sets
	// =========================
	allRoles   := []string{"admin", "manager", "staff"}
	adminOnly  := []string{"admin"}
	adminMgr   := []string{"admin", "manager"}

	// =========================
	// Public Routes (no auth)
	// =========================
	http.HandleFunc("/auth/register", authHandler.Register)
	http.HandleFunc("/auth/login", authHandler.Login)

	// =========================
	// Warehouse Routes
	// =========================
	// GET  — all roles can read
	// POST/PUT/DELETE — admin only
	http.Handle("/warehouses", middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			middleware.RequireRole(allRoles...)(http.HandlerFunc(warehouseHandler.ListWarehouses)).ServeHTTP(w, r)
		case http.MethodPost:
			middleware.RequireRole(adminOnly...)(http.HandlerFunc(warehouseHandler.CreateWarehouse)).ServeHTTP(w, r)
		case http.MethodPut:
			middleware.RequireRole(adminOnly...)(http.HandlerFunc(warehouseHandler.UpdateWarehouse)).ServeHTTP(w, r)
		case http.MethodDelete:
			middleware.RequireRole(adminOnly...)(http.HandlerFunc(warehouseHandler.DeleteWarehouse)).ServeHTTP(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})))

	// =========================
	// User Routes (admin only)
	// =========================
	http.Handle("/users", middleware.Chain(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				userHandler.ListUsers(w, r)
			case http.MethodPost:
				userHandler.CreateUser(w, r)
			default:
				http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			}
		}),
		middleware.RequireRole(adminOnly...),
		middleware.JWTAuth,
	))

	http.Handle("/users/get", middleware.Chain(
		http.HandlerFunc(userHandler.GetUser),
		middleware.RequireRole(adminOnly...),
		middleware.JWTAuth,
	))

	http.Handle("/users/update", middleware.Chain(
		http.HandlerFunc(userHandler.UpdateUser),
		middleware.RequireRole(adminOnly...),
		middleware.JWTAuth,
	))

	http.Handle("/users/delete", middleware.Chain(
		http.HandlerFunc(userHandler.DeleteUser),
		middleware.RequireRole(adminOnly...),
		middleware.JWTAuth,
	))

	// =========================
	// Purchase Routes
	// =========================
	// GET/POST — all roles
	// PUT/DELETE — admin only
	http.Handle("/purchases", middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			middleware.RequireRole(allRoles...)(http.HandlerFunc(purchaseHandler.CreatePurchase)).ServeHTTP(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})))

	// =========================
	// ML Agent Routes
	// =========================
	// admin + manager only
	http.Handle("/ml-agent", middleware.Chain(
		http.HandlerFunc(mlHandler.ProcessQuery),
		middleware.RequireRole(adminMgr...),
		middleware.JWTAuth,
	))

	// =========================
	// Start Server
	// =========================
	port := ":8080"
	log.Printf("🚀 Server starting on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
