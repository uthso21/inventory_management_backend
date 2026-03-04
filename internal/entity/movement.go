package entities

import "time"

type Movement struct {
	ID          int       `json:"id"`
	ProductID   int       `json:"product_id"`
	WarehouseID int       `json:"warehouse_id"`
	UserID      int       `json:"user_id"`
	Type        string    `json:"type"` // PURCHASE or SALE
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`

	// Joined Fields (for frontend display)
	ProductName   string `json:"product_name"`
	WarehouseName string `json:"warehouse_name"`
	UserName      string `json:"user_name"`
}
