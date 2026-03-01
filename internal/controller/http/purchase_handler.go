package http

import (
	"encoding/json"
	"net/http"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
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

func (h *PurchaseHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {
	var purchase entities.Purchase

	if err := json.NewDecoder(r.Body).Decode(&purchase); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	purchaseID, err := h.purchaseService.CreatePurchase(r.Context(), &purchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	purchase.ID = int(purchaseID)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(purchase)
}
