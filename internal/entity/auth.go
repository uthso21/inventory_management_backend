package entities

// LoginRequest is the payload for POST /auth/login
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse is returned on successful login
type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterRequest is the payload for POST /auth/register
type RegisterRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Role        string `json:"role"`         // "admin" | "manager" | "staff"
	WarehouseID *int   `json:"warehouse_id"` // optional for admin
}
