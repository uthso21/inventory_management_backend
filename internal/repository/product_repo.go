package repository

import (
	"context"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
)

type ProductRepository interface {
	ExistsByID(ctx context.Context, id int) (bool, error)
}

type productRepository struct{}

func NewProductRepository() ProductRepository {
	return &productRepository{}
}

func (r *productRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`
	var exists bool
	err := database.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check product existence: %w", err)
	}
	return exists, nil
}
