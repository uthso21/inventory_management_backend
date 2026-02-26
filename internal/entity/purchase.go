package entities

import "time"

type Purchase struct {
	ID          int       `json:"id"`
	WarehouseID int       `json:"warehouse_id"`
	TotalAmount float64   `json:"total_amount"`
	CreatedAt   time.Time `json:"created_at"`
}

type PurchaseItem struct {
	ID         int     `json:"id"`
	PurchaseID int     `json:"purchase_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
}
