package http

import (
	"encoding/json"
	"net/http"

	"github.com/uthso21/inventory_management_backend/internal/service"
)

type MovementHandler struct {
	movementService service.MovementService
}

func NewMovementHandler(service service.MovementService) *MovementHandler {
	return &MovementHandler{
		movementService: service,
	}
}

func (h *MovementHandler) GetMovements(w http.ResponseWriter, r *http.Request) {

	// Query Params
	mType := r.URL.Query().Get("type")
	productID := r.URL.Query().Get("product_id")
	date := r.URL.Query().Get("date")

	movements, err := h.movementService.GetMovements(
		r.Context(),
		mType,
		productID,
		date,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movements)
}
