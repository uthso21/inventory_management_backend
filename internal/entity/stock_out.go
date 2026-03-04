package entities

import "time"

// StockOut entity (records each stock-out transaction)
type StockOut struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	WarehouseID int       `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	Reason      string    `json:"reason"`
	CreatedAt   time.Time `json:"created_at"`
}

// Request DTO for stock-out API
type StockOutRequest struct {
	ProductID   int    `json:"product_id"`
	WarehouseID int    `json:"warehouse_id"`
	Quantity    int    `json:"quantity"`
	Reason      string `json:"reason"`
}