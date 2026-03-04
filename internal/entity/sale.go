package entities

import "time"

type Sale struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	WarehouseID int       `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	SalePrice   float64   `json:"sale_price"`
	CreatedAt   time.Time `json:"created_at"`
}
