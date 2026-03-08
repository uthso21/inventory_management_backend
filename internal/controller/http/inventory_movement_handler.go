package http

import (
	"net/http"
	"strconv"

	"github.com/uthso21/inventory_management_backend/internal/service"
)

// InventoryMovementHandler handles HTTP endpoints for inventory movement history.
// Movements are created internally by PurchaseService (stock-in) and StockOutService (stock-out).
type InventoryMovementHandler struct {
	movementService service.InventoryMovementService
}

func NewInventoryMovementHandler(movementService service.InventoryMovementService) *InventoryMovementHandler {
	return &InventoryMovementHandler{movementService: movementService}
}

// ListMovements handles GET /inventory-movements
// Supports optional query params: product_id, warehouse_id
// Without params → returns all movements (most recent first)
// With product_id  → filters by product
// With warehouse_id → filters by warehouse
func (h *InventoryMovementHandler) ListMovements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ctx := r.Context()

	// Filter by product_id if provided
	if productIDStr := r.URL.Query().Get("product_id"); productIDStr != "" {
		productID, err := strconv.Atoi(productIDStr)
		if err != nil || productID <= 0 {
			writeError(w, http.StatusBadRequest, "invalid product_id")
			return
		}
		movements, err := h.movementService.ListByProduct(ctx, productID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to retrieve movements")
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"movements": movements,
			"total":     len(movements),
		})
		return
	}

	// Filter by warehouse_id if provided
	if warehouseIDStr := r.URL.Query().Get("warehouse_id"); warehouseIDStr != "" {
		warehouseID, err := strconv.Atoi(warehouseIDStr)
		if err != nil || warehouseID <= 0 {
			writeError(w, http.StatusBadRequest, "invalid warehouse_id")
			return
		}
		movements, err := h.movementService.ListByWarehouse(ctx, warehouseID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to retrieve movements")
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"movements": movements,
			"total":     len(movements),
		})
		return
	}

	// No filter — return all
	movements, err := h.movementService.ListAll(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve movements")
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"movements": movements,
		"total":     len(movements),
	})
}
