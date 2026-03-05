package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

type MLHandler struct {
	mlService usecases.MLService
}

func NewMLHandler(mlService usecases.MLService) *MLHandler {
	return &MLHandler{mlService: mlService}
}

func (h *MLHandler) HealthCheck(c echo.Context) error {
	resp, err := h.mlService.HealthCheck(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *MLHandler) DemandForecast(c echo.Context) error {
	var req entities.DemandForecastRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request: " + err.Error()})
	}

	resp, err := h.mlService.GetDemandForecast(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *MLHandler) SmartReorder(c echo.Context) error {
	var req entities.SmartReorderRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request: " + err.Error()})
	}

	resp, err := h.mlService.GetSmartReorder(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *MLHandler) PriceOptimization(c echo.Context) error {
	var req entities.PriceOptimizationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request: " + err.Error()})
	}

	resp, err := h.mlService.GetPriceOptimization(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}
