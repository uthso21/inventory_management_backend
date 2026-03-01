package entities

import "time"

type Purchase struct {
	ID            int       `json:"id"`
	ProductID     int       `json:"product_id"`
	WarehouseID   int       `json:"warehouse_id"`
	Quantity      int       `json:"quantity"`
	PurchasePrice float64   `json:"purchase_price"`
	CreatedAt     time.Time `json:"created_at"`
}

type PurchaseItem struct {
	ID         int     `json:"id"`
	PurchaseID int     `json:"purchase_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}
