package http

import (
	"encoding/json"
	"net/http"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
	"github.com/uthso21/inventory_management_backend/internal/service"
)

// StockOutHandler handles stock out requests
type StockOutHandler struct {
	service *service.StockOutService
}

// Constructor
func NewStockOutHandler(s *service.StockOutService) *StockOutHandler {
	return &StockOutHandler{service: s}
}

// Response structure for frontend
type StockOutResponse struct {
	ProductID        int    `json:"product_id"`
	WarehouseID      int    `json:"warehouse_id"`
	OldQuantity      int    `json:"old_quantity"`
	StockOutQuantity int    `json:"stock_out_quantity"`
	NewQuantity      int    `json:"new_quantity"`
	Message          string `json:"message"`
}

// StockOut handles the POST request for stock out
func (h *StockOutHandler) StockOut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req entities.StockOutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 1️⃣ Get current stock before stock out
	oldQty, err := h.service.GetInventory(req.ProductID, req.WarehouseID)
	if err != nil {
		http.Error(w, "Failed to fetch current inventory: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 2️⃣ Validate quantity
	if req.Quantity <= 0 {
		http.Error(w, "Quantity must be a positive integer", http.StatusBadRequest)
		return
	}
	if req.Quantity > oldQty {
		http.Error(w, "Insufficient stock", http.StatusBadRequest)
		return
	}

	// 3️⃣ Perform stock out
	if err := h.service.StockOutProduct(req.ProductID, req.WarehouseID, req.Quantity, req.Reason); err != nil {
		http.Error(w, "Failed to record stock out: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 4️⃣ Get updated stock after stock out
	newQty, err := h.service.GetInventory(req.ProductID, req.WarehouseID)
	if err != nil {
		http.Error(w, "Stock out succeeded but failed to fetch updated inventory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 5️⃣ Prepare JSON response
	resp := StockOutResponse{
		ProductID:        req.ProductID,
		WarehouseID:      req.WarehouseID,
		OldQuantity:      oldQty,
		StockOutQuantity: req.Quantity,
		NewQuantity:      newQty,
		Message:          "Stock out recorded successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}