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

	// DB init
	database.Connect()

	// Repositories
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	purchaseRepo := repository.NewPurchaseRepository()
	stockOutRepo := repository.NewStockOutRepository() // NEW

	// Services
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, warehouseRepo)
	stockOutService := service.NewStockOutService(stockOutRepo) // NEW

	// Handlers
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)
	stockOutHandler := httpHandler.NewStockOutHandler(stockOutService) // NEW

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/warehouses", warehouseHandler.CreateWarehouse)
	http.HandleFunc("/purchases", purchaseHandler.CreatePurchase)
	http.HandleFunc("/stock-out", stockOutHandler.StockOut)// NEW

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

