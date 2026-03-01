package http

import (
	"encoding/json"
	"errors"
	"net/http"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

type PurchaseHandler struct {
	purchaseService service.PurchaseService
}

type createPurchaseResponse struct {
	PurchaseID int64  `json:"purchase_id"`
	Status     string `json:"status"`
}

func NewPurchaseHandler(purchaseService service.PurchaseService) *PurchaseHandler {
	return &PurchaseHandler{
		purchaseService: purchaseService,
	}
}

func (h *PurchaseHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var purchase entities.Purchase

	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	purchaseID, err := h.purchaseService.CreatePurchase(r.Context(), &purchase)
	if err != nil {
		// service validation or DB constraint fail -> 400
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	purchase.ID = int(purchaseID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(purchase)
}
