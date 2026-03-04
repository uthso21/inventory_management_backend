package http

import (
	"encoding/json"
	"net/http"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

// MLAgentHandler handles HTTP requests for ML agent operations
type MLAgentHandler struct {
	mlService usecases.MLAgentService
}

// NewMLAgentHandler creates a new instance of MLAgentHandler
func NewMLAgentHandler(mlService usecases.MLAgentService) *MLAgentHandler {
	return &MLAgentHandler{
		mlService: mlService,
	}
}

// ProcessQuery handles POST /ml/agent
// This is the main endpoint that receives data from frontend,
// forwards it to the FastAPI microservice, and returns the result
func (h *MLAgentHandler) ProcessQuery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req entities.MLAgentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Validate required fields
	if req.Query == "" {
		SendError(w, http.StatusBadRequest, "query is required")
		return
	}
	if req.Context.ProductID == "" {
		SendError(w, http.StatusBadRequest, "context.product_id is required")
		return
	}

	// Forward to ML microservice
	resp, err := h.mlService.ProcessQuery(r.Context(), &req)
	if err != nil {
		SendError(w, http.StatusServiceUnavailable, "ML service error: "+err.Error())
		return
	}

	SendSuccess(w, "Analysis complete", resp)
}

// DemandForecast handles POST /ml/demand-forecast
// Convenience endpoint for demand forecasting
func (h *MLAgentHandler) DemandForecast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var productCtx entities.ProductContext
	if err := json.NewDecoder(r.Body).Decode(&productCtx); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if productCtx.ProductID == "" {
		SendError(w, http.StatusBadRequest, "product_id is required")
		return
	}

	resp, err := h.mlService.GetDemandForecast(r.Context(), &productCtx)
	if err != nil {
		SendError(w, http.StatusServiceUnavailable, "ML service error: "+err.Error())
		return
	}

	SendSuccess(w, "Demand forecast complete", resp)
}

// SmartReorder handles POST /ml/smart-reorder
// Convenience endpoint for smart reorder recommendations
func (h *MLAgentHandler) SmartReorder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var productCtx entities.ProductContext
	if err := json.NewDecoder(r.Body).Decode(&productCtx); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if productCtx.ProductID == "" {
		SendError(w, http.StatusBadRequest, "product_id is required")
		return
	}

	resp, err := h.mlService.GetSmartReorder(r.Context(), &productCtx)
	if err != nil {
		SendError(w, http.StatusServiceUnavailable, "ML service error: "+err.Error())
		return
	}

	SendSuccess(w, "Smart reorder analysis complete", resp)
}

// PricelistOptimize handles POST /ml/pricelist-optimize
// Convenience endpoint for pricelist optimization
func (h *MLAgentHandler) PricelistOptimize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var productCtx entities.ProductContext
	if err := json.NewDecoder(r.Body).Decode(&productCtx); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if productCtx.ProductID == "" {
		SendError(w, http.StatusBadRequest, "product_id is required")
		return
	}

	resp, err := h.mlService.GetPricelistOptimization(r.Context(), &productCtx)
	if err != nil {
		SendError(w, http.StatusServiceUnavailable, "ML service error: "+err.Error())
		return
	}

	SendSuccess(w, "Pricelist optimization complete", resp)
}

// FullAnalysis handles POST /ml/full-analysis
// Runs all three ML tools
func (h *MLAgentHandler) FullAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var productCtx entities.ProductContext
	if err := json.NewDecoder(r.Body).Decode(&productCtx); err != nil {
		SendError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	if productCtx.ProductID == "" {
		SendError(w, http.StatusBadRequest, "product_id is required")
		return
	}

	resp, err := h.mlService.GetFullAnalysis(r.Context(), &productCtx)
	if err != nil {
		SendError(w, http.StatusServiceUnavailable, "ML service error: "+err.Error())
		return
	}

	SendSuccess(w, "Full analysis complete", resp)
}

// HealthCheck handles GET /ml/health
// Checks if the ML microservice is available
func (h *MLAgentHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	isHealthy, err := h.mlService.HealthCheck(r.Context())
	if err != nil || !isHealthy {
		SendJSON(w, http.StatusServiceUnavailable, Response{
			Success: false,
			Error:   "ML service is unavailable",
			Data: map[string]interface{}{
				"ml_service": "down",
				"error":      err.Error(),
			},
		})
		return
	}

	SendSuccess(w, "ML service is healthy", map[string]interface{}{
		"ml_service": "up",
		"go_backend": "up",
	})
}
