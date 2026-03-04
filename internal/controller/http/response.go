package http

import (
	"encoding/json"
	"net/http"
)

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SendJSON sends a JSON response
func SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// SendSuccess sends a success response
func SendSuccess(w http.ResponseWriter, message string, data interface{}) {
	SendJSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SendError sends an error response
func SendError(w http.ResponseWriter, statusCode int, message string) {
	SendJSON(w, statusCode, Response{
		Success: false,
		Error:   message,
	})
}

// writeError writes a compact JSON error (used internally by handlers).
func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// writeJSON writes any value as JSON (used internally by handlers).
func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
