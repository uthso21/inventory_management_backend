package main

import (
	"log"
	"net/http"

	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {

	// =========================
	// User Module Initialization
	// =========================

	userRepo := repository.NewUserRepository()
	userService := service.NewUserService(userRepo)
	userHandler := httpHandler.NewUserHandler(userService)

	// =========================
	// Purchase Module Initialization
	// =========================

	purchaseRepo := repository.NewPurchaseRepository()
	purchaseService := service.NewPurchaseService(purchaseRepo)
	purchaseHandler := httpHandler.NewPurchaseHandler(purchaseService)

	// =========================
	// Routes
	// =========================

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

	// Purchase route
	http.HandleFunc("/purchases", purchaseHandler.CreatePurchase)

	// =========================
	// Start Server
	// =========================

	port := ":8080"
	log.Printf("Server starting on port %s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
