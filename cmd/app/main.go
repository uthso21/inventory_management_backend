package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/uthso21/inventory_management_backend/docs"
	httpHandler "github.com/uthso21/inventory_management_backend/internal/controller/http"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

// @title Inventory Management ML API
// @version 1.0
// @description API for ML-powered inventory management features including demand forecasting, smart reorder, and price optimization.

// @host localhost:8080
// @BasePath /

func main() {
	mlService := usecases.NewMLService()
	mlHandler := httpHandler.NewMLHandler(mlService)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Swagger UI
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// ML Routes
	e.GET("/ml/health", mlHandler.HealthCheck)
	e.POST("/ml/demand-forecast", mlHandler.DemandForecast)
	e.POST("/ml/smart-reorder", mlHandler.SmartReorder)
	e.POST("/ml/price-optimization", mlHandler.PriceOptimization)

	log.Println("Server starting on :8080")
	log.Println("Swagger UI available at http://localhost:8080/swagger/index.html")
	e.Logger.Fatal(e.Start(":8080"))
}
