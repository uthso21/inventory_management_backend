package routes

import (
	"net/http"

	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/middleware"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

// Setup wires all dependencies and registers HTTP routes.
func Setup() {
	// Repositories
	userRepo := repository.NewUserRepository()
	warehouseRepo := repository.NewWarehouseRepository()
	purchaseRepo := repository.NewPurchaseRepository()
	stockOutRepo := repository.NewStockOutRepository()
	productRepo := repository.NewProductRepository()
	inventoryMovementRepo := repository.NewInventoryMovementRepository()

	// Services
	userService := service.NewUserService(userRepo)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	purchaseService := service.NewPurchaseService(purchaseRepo, warehouseRepo, productRepo, inventoryMovementRepo)
	productService := service.NewProductService(productRepo)
	stockOutService := service.NewStockOutService(stockOutRepo)
	mlService := service.NewMLAgentServiceWithDefaults()

	// Handlers
	authHandler := httpHandler.NewAuthHandler(userService)
	userHandler := httpHandler.NewUserHandler(userService)
	warehouseHandler := httpHandler.NewWarehouseHandler(warehouseService)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)
	stockOutHandler := httpHandler.NewStockOutHandler(stockOutService)
	productHandler := httpHandler.NewProductHandler(productService)
	mlHandler := httpHandler.NewMLAgentHandler(mlService)

	// Auth routes
	http.HandleFunc("/auth/register", authHandler.Register)
	http.HandleFunc("/auth/login", authHandler.Login)

	// User routes
	http.HandleFunc("/users", userHandler.CreateUser)

	// Warehouse routes
	http.HandleFunc("/warehouses", warehouseHandler.CreateWarehouse)

	// Purchase routes (JWT protected)
	http.Handle("/purchases", middleware.JWTAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			purchaseHandler.ListPurchases(w, r)
		case http.MethodPost:
			purchaseHandler.CreatePurchase(w, r)
		default:
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
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
			http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/products/get", productHandler.GetProduct)
	http.HandleFunc("/products/update", productHandler.UpdateProduct)
	http.HandleFunc("/products/delete", productHandler.DeleteProduct)
	http.HandleFunc("/products/low-stock", productHandler.GetLowStockProducts)

	// Stock-out routes
	http.HandleFunc("/stock-out", stockOutHandler.StockOut)

	// ML agent routes
	http.HandleFunc("/ml/agent", mlHandler.ProcessQuery)
	http.HandleFunc("/ml/demand-forecast", mlHandler.DemandForecast)
	http.HandleFunc("/ml/smart-reorder", mlHandler.SmartReorder)
	http.HandleFunc("/ml/pricelist-optimize", mlHandler.PricelistOptimize)
	http.HandleFunc("/ml/full-analysis", mlHandler.FullAnalysis)
	http.HandleFunc("/ml/health", mlHandler.HealthCheck)
}
