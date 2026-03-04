package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

type MLAgentHandler struct {
	mlService usecases.MLAgentService
}

func NewMLAgentHandler(mlService usecases.MLAgentService) *MLAgentHandler {
	return &MLAgentHandler{mlService: mlService}
}

func (h *MLAgentHandler) HealthCheck(c echo.Context) error {
resp, err := h.mlService.HealthCheck(c.Request().Context())
if err != nil {
return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": err.Error()})
}
return c.JSON(http.StatusOK, resp)
}

func (h *MLAgentHandler) DemandForecast(c echo.Context) error {
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

func (h *MLAgentHandler) SmartReorder(c echo.Context) error {
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

func (h *MLAgentHandler) PriceOptimization(c echo.Context) error {
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
