package main

import (
	"log"
	"net/http"

	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {
	// Initialize repository layer
	userRepo := repository.NewUserRepository()
	productRepo := repository.NewProductRepository()

	// Initialize use case/service layer
	userService := usecases.NewUserService(userRepo)
	productService := usecases.NewProductService(productRepo)

	// Initialize HTTP handler layer
	userHandler := httpHandler.NewUserHandler(userService)
	productHandler := httpHandler.NewProductHandler(productService)

	// User routes
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

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
