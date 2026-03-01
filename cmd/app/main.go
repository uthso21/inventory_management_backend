package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/database"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {
	_ = godotenv.Load()

	// DB init
	database.Connect()

	// Repositories
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	purchaseRepo := repository.NewPurchaseRepository()
	productRepo := repository.NewProductRepository()

	// Services
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, warehouseRepo, productRepo)

	// Handlers
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/warehouses", warehouseHandler.CreateWarehouse)
	http.HandleFunc("/purchases", purchaseHandler.CreatePurchase)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
