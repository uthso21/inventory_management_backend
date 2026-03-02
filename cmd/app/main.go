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

	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DB init
	database.Connect()

	// Repositories
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	purchaseRepo := repository.NewPurchaseRepository()
	productRepo := repository.NewProductRepository()
	inventoryMovementRepo := repository.NewInventoryMovementRepository()

	// Services
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, warehouseRepo, productRepo, inventoryMovementRepo)
	productService := service.NewProductService(productRepo)

	// Handlers
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)
	productHandler := httpHandler.NewProductHandler(productService)

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/warehouses", warehouseHandler.CreateWarehouse)

	// Purchase routes with authentication
	http.Handle("/purchases", middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			purchaseHandler.ListPurchases(w, r)
		case http.MethodPost:
			purchaseHandler.CreatePurchase(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))
	http.Handle("/purchases/get", middleware.JWTAuth(http.HandlerFunc(purchaseHandler.GetPurchase)))

	// Product routes
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			productHandler.ListProducts(w, r)
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/products/get", productHandler.GetProduct)
	http.HandleFunc("/products/update", productHandler.UpdateProduct)
	http.HandleFunc("/products/delete", productHandler.DeleteProduct)
	http.HandleFunc("/products/low-stock", productHandler.GetLowStockProducts)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
