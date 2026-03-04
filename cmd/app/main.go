package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/uthso21/inventory_management_backend/internal/database"
	"github.com/uthso21/inventory_management_backend/internal/routes"
)

func main() {
	// Load env if a .env file is present; avoid noisy log when it's absent
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Error loading .env: %v", err)
		}
	}

	// DB init
	database.Connect()

	// Register all routes
	routes.Setup()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

