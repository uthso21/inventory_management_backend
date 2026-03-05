package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

type MLHandler struct {
	mlService usecases.MLAgentService
}

func NewMLHandler(mlService usecases.MLAgentService) *MLHandler {
	return &MLHandler{mlService: mlService}
}

// @Summary Check ML service health
// @Tags ML
// @Produce json
// @Success 200 {object} entities.MLHealthResponse
// @Failure 503 {object} map[string]string
// @Router /ml/health [get]
func (h *MLHandler) HealthCheck(c echo.Context) error {
	resp, err := h.mlService.HealthCheck(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, resp)
}

// @Summary Predict demand for a product
// @Tags ML
// @Accept json
// @Produce json
// @Param request body entities.DemandForecastRequest true "Demand Forecast Request"
// @Success 200 {object} entities.DemandForecastResponse
// @Failure 400 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /ml/demand-forecast [post]
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

// @Summary Calculate smart reorder quantity
// @Tags ML
// @Accept json
// @Produce json
// @Param request body entities.SmartReorderRequest true "Smart Reorder Request"
// @Success 200 {object} entities.SmartReorderResponse
// @Failure 400 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /ml/smart-reorder [post]
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

// @Summary Optimize product pricing
// @Tags ML
// @Accept json
// @Produce json
// @Param request body entities.PriceOptimizationRequest true "Price Optimization Request"
// @Success 200 {object} entities.PriceOptimizationResponse
// @Failure 400 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /ml/price-optimization [post]
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
