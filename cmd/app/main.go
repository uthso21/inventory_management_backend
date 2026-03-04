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

	// Load .env (optional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect Database
	database.Connect()

	// ---------------- Repositories ----------------
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	productRepo := repository.NewProductRepository()
	purchaseRepo := repository.NewPurchaseRepository()
	movementRepo := repository.NewMovementRepository()

	// ---------------- Services ----------------
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	productService := service.NewProductService(productRepo)

	purchaseService := service.NewPurchaseService(
		purchaseRepo,
		warehouseRepo,
		productRepo,
		movementRepo,
	)

	movementService := service.NewMovementService(movementRepo)

	// ---------------- Handlers ----------------
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	productHandler := httpHandler.NewProductHandler(productService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)
	movementHandler := httpHandler.NewMovementHandler(movementService)

	// ---------------- Routes ----------------

	// Users
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			userHandler.CreateUser(w, r)
		case http.MethodGet:
			userHandler.ListUsers(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Warehouses
	http.HandleFunc("/warehouses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			warehouseHandler.CreateWarehouse(w, r)
		case http.MethodGet:
			warehouseHandler.ListWarehouses(w, r)
		case http.MethodPut:
			warehouseHandler.UpdateWarehouse(w, r)
		case http.MethodDelete:
			warehouseHandler.DeleteWarehouse(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Products
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			productHandler.CreateProduct(w, r)
		case http.MethodGet:
			productHandler.ListProducts(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/products/get", productHandler.GetProduct)
	http.HandleFunc("/products/update", productHandler.UpdateProduct)
	http.HandleFunc("/products/delete", productHandler.DeleteProduct)
	http.HandleFunc("/products/low-stock", productHandler.GetLowStockProducts)

	// Purchases
	http.HandleFunc("/purchases", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			purchaseHandler.CreatePurchase(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Movements (NEW)
	http.HandleFunc("/movements", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			movementHandler.GetMovements(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// ---------------- Start Server ----------------
	log.Println(" Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
