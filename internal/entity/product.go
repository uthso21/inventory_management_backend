package entities

import "time"

// Product represents a product in the inventory
type Product struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	SKU          string    `json:"sku" db:"sku"`
	Description  string    `json:"description" db:"description"`
	Category     string    `json:"category" db:"category"`
	UnitPrice    float64   `json:"unit_price" db:"unit_price"`
	CurrentStock int       `json:"current_stock" db:"current_stock"`
	ReorderPoint int       `json:"reorder_point" db:"reorder_point"`
	LeadTimeDays int       `json:"lead_time_days" db:"lead_time_days"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// CreateProductRequest for creating a new product
type CreateProductRequest struct {
	Name         string  `json:"name" validate:"required"`
	SKU          string  `json:"sku" validate:"required"`
	Description  string  `json:"description"`
	Category     string  `json:"category"`
	UnitPrice    float64 `json:"unit_price" validate:"required,gt=0"`
	CurrentStock int     `json:"current_stock" validate:"gte=0"`
	ReorderPoint int     `json:"reorder_point" validate:"gte=0"`
	LeadTimeDays int     `json:"lead_time_days" validate:"gte=1"`
}

// UpdateProductRequest for updating a product
type UpdateProductRequest struct {
	Name         *string  `json:"name,omitempty"`
	Description  *string  `json:"description,omitempty"`
	Category     *string  `json:"category,omitempty"`
	UnitPrice    *float64 `json:"unit_price,omitempty"`
	CurrentStock *int     `json:"current_stock,omitempty"`
	ReorderPoint *int     `json:"reorder_point,omitempty"`
	LeadTimeDays *int     `json:"lead_time_days,omitempty"`
}

// SalesHistory represents historical sales data
type SalesHistory struct {
	ID           int       `json:"id" db:"id"`
	ProductID    int       `json:"product_id" db:"product_id"`
	WeekNumber   int       `json:"week_number" db:"week_number"`
	Year         int       `json:"year" db:"year"`
	QuantitySold int       `json:"quantity_sold" db:"quantity_sold"`
	Revenue      float64   `json:"revenue" db:"revenue"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// InventoryTransaction represents stock movement
type InventoryTransaction struct {
	ID              int       `json:"id" db:"id"`
	ProductID       int       `json:"product_id" db:"product_id"`
	TransactionType string    `json:"transaction_type" db:"transaction_type"`
	Quantity        int       `json:"quantity" db:"quantity"`
	Notes           string    `json:"notes" db:"notes"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
