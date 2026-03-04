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

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB connect
	database.Connect()

	// Repositories
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	productRepo := repository.NewProductRepository()
	purchaseRepo := repository.NewPurchaseRepository()
	inventoryMovementRepo := repository.NewInventoryMovementRepository()

	// Services
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	productService := service.NewProductService(productRepo)

	purchaseService := service.NewPurchaseService(
		purchaseRepo,
		warehouseRepo,
		productRepo,
		inventoryMovementRepo,
	)

	// Handlers
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	productHandler := httpHandler.NewProductHandler(productService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/warehouses", warehouseHandler.CreateWarehouse)
	http.HandleFunc("/products", productHandler.CreateProduct)
	http.HandleFunc("/purchases", purchaseHandler.CreatePurchase)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
