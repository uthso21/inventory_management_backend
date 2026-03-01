package entities

import "time"

// User represents the user entity in the system
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`            // bcrypt hash — never serialized
	Role         string    `json:"role"`         // "admin" | "manager" | "staff"
	WarehouseID  *int      `json:"warehouse_id"` // nullable (admin has no warehouse restriction)
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
