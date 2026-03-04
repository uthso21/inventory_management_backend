package entities

import "time"

// Purchase represents a stock-in transaction
type Purchase struct {
	ID          int            `json:"id"`
	WarehouseID int            `json:"warehouse_id"`
	CreatedBy   int            `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Items       []PurchaseItem `json:"items,omitempty"`
}

// PurchaseItem represents a line item in a purchase
type PurchaseItem struct {
	ID         int      `json:"id"`
	PurchaseID int      `json:"purchase_id"`
	ProductID  int      `json:"product_id"`
	Quantity   int      `json:"quantity"`
	UnitPrice  *float64 `json:"unit_price,omitempty"`
}

// CreatePurchaseRequest represents the request body for creating a purchase
type CreatePurchaseRequest struct {
	WarehouseID int            `json:"warehouse_id"`
	Items       []PurchaseItem `json:"items"`
}

// InventoryMovement represents a stock movement log entry
type InventoryMovement struct {
	ID            int       `json:"id"`
	ProductID     int       `json:"product_id"`
	WarehouseID   int       `json:"warehouse_id"`
	MovementType  string    `json:"movement_type"` // purchase, sale, adjustment, transfer
	Quantity      int       `json:"quantity"`
	ReferenceType string    `json:"reference_type,omitempty"`
	ReferenceID   int       `json:"reference_id,omitempty"`
	CreatedBy     int       `json:"created_by"`
	Notes         string    `json:"notes,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
