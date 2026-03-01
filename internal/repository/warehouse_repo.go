package repository

import (
	"context"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *entities.Warehouse) error
	List(ctx context.Context) ([]*entities.Warehouse, error)
	Update(ctx context.Context, warehouse *entities.Warehouse) error
	Delete(ctx context.Context, id int) error
	ExistsByID(ctx context.Context, id int) (bool, error)
}

type warehouseRepository struct{}

func NewWarehouseRepository() WarehouseRepository {
	return &warehouseRepository{}
}

// -------------------- CREATE --------------------

func (r *warehouseRepository) Create(ctx context.Context, warehouse *entities.Warehouse) error {

	query := `
		INSERT INTO warehouses (name, location, description)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`

	err := database.DB.QueryRowContext(
		ctx,
		query,
		warehouse.Name,
		warehouse.Location,
		warehouse.Description,
	).Scan(&warehouse.ID, &warehouse.CreatedAt, &warehouse.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create warehouse: %w", err)
	}

	return nil
}

// -------------------- LIST --------------------

func (r *warehouseRepository) List(ctx context.Context) ([]*entities.Warehouse, error) {

	query := `
		SELECT id, name, location, description, created_at, updated_at
		FROM warehouses
		ORDER BY id DESC
	`

	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var warehouses []*entities.Warehouse

	for rows.Next() {
		var warehouse entities.Warehouse
		err := rows.Scan(
			&warehouse.ID,
			&warehouse.Name,
			&warehouse.Location,
			&warehouse.Description,
			&warehouse.CreatedAt,
			&warehouse.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, &warehouse)
	}

	return warehouses, nil
}

// -------------------- UPDATE --------------------

func (r *warehouseRepository) Update(ctx context.Context, warehouse *entities.Warehouse) error {

	query := `
		UPDATE warehouses
		SET name = $1,
		    location = $2,
		    description = $3,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
		RETURNING created_at, updated_at
	`

	err := database.DB.QueryRowContext(
		ctx,
		query,
		warehouse.Name,
		warehouse.Location,
		warehouse.Description,
		warehouse.ID,
	).Scan(&warehouse.CreatedAt, &warehouse.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update warehouse: %w", err)
	}

	return nil
}

// -------------------- DELETE --------------------

func (r *warehouseRepository) Delete(ctx context.Context, id int) error {

	query := `DELETE FROM warehouses WHERE id = $1`

	result, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("warehouse not found")
	}

	return nil
}

// -------------------- EXISTS --------------------

func (r *warehouseRepository) ExistsByID(ctx context.Context, id int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM warehouses WHERE id = $1)`
	var exists bool
	err := database.DB.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check warehouse existence: %w", err)
	}
	return exists, nil
}
