package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	usecases "github.com/uthso21/inventory_management_backend/internal/service"
)

type WarehouseHandler struct {
	service usecases.WarehouseService
}

func NewWarehouseHandler(service usecases.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

// -------------------- CREATE --------------------

func (h *WarehouseHandler) CreateWarehouse(w http.ResponseWriter, r *http.Request) {
	var warehouse entities.Warehouse

	err := json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.CreateWarehouse(r.Context(), &warehouse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(warehouse)
}

// -------------------- LIST --------------------

func (h *WarehouseHandler) ListWarehouses(w http.ResponseWriter, r *http.Request) {

	warehouses, err := h.service.ListWarehouses(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(warehouses)
}

// -------------------- UPDATE --------------------

func (h *WarehouseHandler) UpdateWarehouse(w http.ResponseWriter, r *http.Request) {

	var warehouse entities.Warehouse

	err := json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = h.service.UpdateWarehouse(r.Context(), &warehouse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(warehouse)
}

// -------------------- DELETE --------------------

func (h *WarehouseHandler) DeleteWarehouse(w http.ResponseWriter, r *http.Request) {

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteWarehouse(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
