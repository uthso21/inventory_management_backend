package entities

import "time"

type StockOut struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	WarehouseID int       `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	Reason      string    `json:"reason"` // optional (sale, damage, etc.)
	CreatedAt   time.Time `json:"created_at"`
}