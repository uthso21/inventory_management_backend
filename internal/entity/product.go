package entities

type Product struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	SKU          string  `json:"sku"`
	Price        float64 `json:"price"`
	Description  string  `json:"description"`
	Stock        int     `json:"stock"`
	ReorderLevel int     `json:"reorder_level"`
}
