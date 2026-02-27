package entities

import "time"

// Warehouse represents the warehouse entity in the system
type Warehouse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
