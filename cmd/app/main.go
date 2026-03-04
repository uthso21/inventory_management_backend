package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/uthso21/inventory_management_backend/configs"
	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	"github.com/uthso21/inventory_management_backend/internal/repository"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

func main() {
	dbConfig := configs.LoadDatabaseConfig()
	db, err := configs.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	_ = productRepo

	userService := usecases.NewUserService(userRepo)
	mlService := usecases.NewMLAgentServiceWithDefaults()

	userHandler := httpHandler.NewUserHandler(userService)
	mlHandler := httpHandler.NewMLAgentHandler(mlService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/users", userHandler.ListUsers)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	e.GET("/ml/health", mlHandler.HealthCheck)
	e.POST("/ml/demand-forecast", mlHandler.DemandForecast)
	e.POST("/ml/smart-reorder", mlHandler.SmartReorder)
	e.POST("/ml/price-optimization", mlHandler.PriceOptimization)

	log.Println("Server starting on :8080")
	e.Logger.Fatal(e.Start(":8080"))
}
