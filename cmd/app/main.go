package main

import (
	"log"
	"net/http"

	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/database"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {

	// =========================
	// Connect to Database FIRST
	// =========================
	database.Connect()

	// =========================
	// Warehouse Setup
	// =========================
	warehouseRepo := repository.NewWarehouseRepository()
	warehouseService := service.NewWarehouseService(warehouseRepo)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)

	// =========================
	// User Setup
	// =========================
	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := httpHandler.NewUserHandler(userService)

	// =========================
	// Purchase Setup
	// =========================
	purchaseRepo := repository.NewPurchaseRepository()
	purchaseService := service.NewPurchaseService(purchaseRepo)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)

	// =========================
	// Routes
	// =========================

	// Users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.ListUsers(w, r)
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/get", userHandler.GetUser)
	http.HandleFunc("/users/update", userHandler.UpdateUser)
	http.HandleFunc("/users/delete", userHandler.DeleteUser)

	// Purchases
	http.HandleFunc("/purchases", purchaseHandler.CreatePurchase)

	// Warehouses
	http.HandleFunc("/warehouses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			warehouseHandler.ListWarehouses(w, r)
		case http.MethodPost:
			warehouseHandler.CreateWarehouse(w, r)
		case http.MethodPut:
			warehouseHandler.UpdateWarehouse(w, r)
		case http.MethodDelete:
			warehouseHandler.DeleteWarehouse(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// =========================
	// Start Server
	// =========================
	port := ":8080"
	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
