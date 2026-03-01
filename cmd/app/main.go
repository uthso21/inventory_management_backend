package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/database"
<<<<<<< HEAD
=======
	"github.com/uthso21/inventory_management_backend/internal/middleware"
>>>>>>> 788bf3e4c93ab81a87fe612c22686b7dffef2020
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {
<<<<<<< HEAD
	// Load env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

=======
>>>>>>> 788bf3e4c93ab81a87fe612c22686b7dffef2020
	// DB init
	database.Connect()

	// Repositories
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	purchaseRepo := repository.NewPurchaseRepository()
<<<<<<< HEAD
	productRepo := repository.NewProductRepository()
=======
>>>>>>> 788bf3e4c93ab81a87fe612c22686b7dffef2020

	// Services
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, warehouseRepo)
<<<<<<< HEAD
	productService := service.NewProductService(productRepo)
=======
>>>>>>> 788bf3e4c93ab81a87fe612c22686b7dffef2020

	// Handlers
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)
<<<<<<< HEAD
	productHandler := httpHandler.NewProductHandler(productService)
=======
>>>>>>> 788bf3e4c93ab81a87fe612c22686b7dffef2020

	// Routes
	http.HandleFunc("/users", userHandler.CreateUser)
	http.HandleFunc("/warehouses", warehouseHandler.CreateWarehouse)
	http.HandleFunc("/purchases", purchaseHandler.CreatePurchase)
<<<<<<< HEAD

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
=======
>>>>>>> 788bf3e4c93ab81a87fe612c22686b7dffef2020

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
