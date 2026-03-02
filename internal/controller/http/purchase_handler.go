package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/middleware"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

type PurchaseHandler struct {
	purchaseService service.PurchaseService
}

func NewPurchaseHandler(purchaseService service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}

// CreatePurchase handles POST /purchases
// Implements task #40: Create purchase API
func (h *PurchaseHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var req entities.CreatePurchaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.WarehouseID == 0 {
		http.Error(w, `{"error":"warehouse_id is required"}`, http.StatusBadRequest)
		return
	}

	if len(req.Items) == 0 {
		http.Error(w, `{"error":"items are required"}`, http.StatusBadRequest)
		return
	}

	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(int)
	if !ok || userID == 0 {
		http.Error(w, `{"error":"unauthorized: user_id not found"}`, http.StatusUnauthorized)
		return
	}

	purchase, err := h.purchaseService.CreatePurchase(r.Context(), &req, userID)
	if err != nil {
		// Determine appropriate error code
		switch {
		case err == service.ErrWarehouseNotFound:
			http.Error(w, `{"error":"warehouse not found"}`, http.StatusNotFound)
		case err == service.ErrProductNotFound || err.Error()[:17] == "product not found":
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusNotFound)
		case err == service.ErrInvalidQuantity:
			http.Error(w, `{"error":"quantity must be greater than zero"}`, http.StatusBadRequest)
		case err == service.ErrEmptyPurchaseItems:
			http.Error(w, `{"error":"purchase items are required"}`, http.StatusBadRequest)
		default:
			http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "purchase created successfully",
		"purchase": purchase,
	})
}

// ListPurchases handles GET /purchases
// Implements task #48: Display purchase history
func (h *PurchaseHandler) ListPurchases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	purchases, err := h.purchaseService.ListPurchases(r.Context())
	if err != nil {
		http.Error(w, `{"error":"failed to retrieve purchases"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"purchases": purchases,
		"total":     len(purchases),
	})
}

// GetPurchase handles GET /purchases/{id}
func (h *PurchaseHandler) GetPurchase(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from query parameter
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"id is required"}`, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	purchase, err := h.purchaseService.GetPurchase(r.Context(), id)
	if err != nil {
		http.Error(w, `{"error":"purchase not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(purchase)
}
