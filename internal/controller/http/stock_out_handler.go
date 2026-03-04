package http

import (
    "encoding/json"
    "net/http"

    "github.com/uthso21/inventory_management_backend/internal/service"
)

type StockOutHandler struct {
    service *service.StockOutService
}

func NewStockOutHandler(s *service.StockOutService) *StockOutHandler {
    return &StockOutHandler{service: s}
}

type StockOutRequest struct {
    ProductID   int    `json:"product_id"`
    WarehouseID int    `json:"warehouse_id"`
    Quantity    int    `json:"quantity"`
    Reason      string `json:"reason"`
}

func (h *StockOutHandler) StockOut(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req StockOutRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    err := h.service.StockOutProduct(req.ProductID, req.WarehouseID, req.Quantity, req.Reason)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Stock out recorded successfully"))
}