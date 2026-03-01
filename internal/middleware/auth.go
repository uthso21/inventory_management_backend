package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// contextKey is an unexported type to avoid collisions in context values
type contextKey string

const (
	ContextKeyUserID      contextKey = "user_id"
	ContextKeyRole        contextKey = "role"
	ContextKeyWarehouseID contextKey = "warehouse_id"
)

// JWTAuth validates the Bearer token and loads claims into the request context.
// Downstream handlers can read user_id, role, warehouse_id via the typed context keys.
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"authorization header required"}`, http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, `{"error":"invalid token claims"}`, http.StatusUnauthorized)
			return
		}

		// Store each claim as a typed context value
		ctx := r.Context()
		if v, ok := claims["user_id"].(float64); ok {
			ctx = context.WithValue(ctx, ContextKeyUserID, int(v))
		}
		if v, ok := claims["role"].(string); ok {
			ctx = context.WithValue(ctx, ContextKeyRole, v)
		}
		if v, ok := claims["warehouse_id"].(float64); ok {
			wid := int(v)
			ctx = context.WithValue(ctx, ContextKeyWarehouseID, &wid)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole returns a middleware that allows only the specified roles.
// Must be used after JWTAuth (which populates the role in context).
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value(ContextKeyRole).(string)
			if !ok || role == "" {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			if _, permitted := allowed[role]; !permitted {
				http.Error(w, `{"error":"forbidden: insufficient role"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// Chain applies a list of middlewares to a handler in order.
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
