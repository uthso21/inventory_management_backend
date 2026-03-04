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

func (h *PurchaseHandler) CreatePurchase(w http.ResponseWriter, r *http.Request) {

	var req entities.CreatePurchaseRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.ContextKeyUserID).(int)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	purchase, err := h.purchaseService.CreatePurchase(r.Context(), &req, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(purchase)
}

func (h *PurchaseHandler) ListPurchases(w http.ResponseWriter, r *http.Request) {

	purchases, err := h.purchaseService.ListPurchases(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch purchases", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(purchases)
}

func (h *PurchaseHandler) GetPurchase(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	purchase, err := h.purchaseService.GetPurchase(r.Context(), id)
	if err != nil {
		http.Error(w, "purchase not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(purchase)
}
