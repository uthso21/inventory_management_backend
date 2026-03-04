package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type InventoryMovementRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, movement *entities.InventoryMovement) error
	List(ctx context.Context) ([]*entities.InventoryMovement, error)
	ListByProductID(ctx context.Context, productID int) ([]*entities.InventoryMovement, error)
	ListByWarehouseID(ctx context.Context, warehouseID int) ([]*entities.InventoryMovement, error)
}

type inventoryMovementRepository struct{}

func NewInventoryMovementRepository() InventoryMovementRepository {
	return &inventoryMovementRepository{}
}

func (r *inventoryMovementRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, movement *entities.InventoryMovement) error {
	query := `
		INSERT INTO inventory_movements (product_id, warehouse_id, movement_type, quantity, reference_type, reference_id, created_by, notes, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, created_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		movement.ProductID,
		movement.WarehouseID,
		movement.MovementType,
		movement.Quantity,
		nullableString(movement.ReferenceType),
		nullableInt(movement.ReferenceID),
		movement.CreatedBy,
		nullableString(movement.Notes),
	).Scan(&movement.ID, &movement.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create inventory movement: %w", err)
	}

	return nil
}

func (r *inventoryMovementRepository) List(ctx context.Context) ([]*entities.InventoryMovement, error) {
	query := `
		SELECT id, product_id, warehouse_id, movement_type, quantity, 
		       COALESCE(reference_type, '') as reference_type, 
		       COALESCE(reference_id, 0) as reference_id, 
		       created_by, COALESCE(notes, '') as notes, created_at
		FROM inventory_movements
		ORDER BY created_at DESC
	`

	return r.queryMovements(ctx, query)
}

func (r *inventoryMovementRepository) ListByProductID(ctx context.Context, productID int) ([]*entities.InventoryMovement, error) {
	query := `
		SELECT id, product_id, warehouse_id, movement_type, quantity, 
		       COALESCE(reference_type, '') as reference_type, 
		       COALESCE(reference_id, 0) as reference_id, 
		       created_by, COALESCE(notes, '') as notes, created_at
		FROM inventory_movements
		WHERE product_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.DB.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory movements: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

func (r *inventoryMovementRepository) ListByWarehouseID(ctx context.Context, warehouseID int) ([]*entities.InventoryMovement, error) {
	query := `
		SELECT id, product_id, warehouse_id, movement_type, quantity, 
		       COALESCE(reference_type, '') as reference_type, 
		       COALESCE(reference_id, 0) as reference_id, 
		       created_by, COALESCE(notes, '') as notes, created_at
		FROM inventory_movements
		WHERE warehouse_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.DB.QueryContext(ctx, query, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory movements: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

func (r *inventoryMovementRepository) queryMovements(ctx context.Context, query string) ([]*entities.InventoryMovement, error) {
	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list inventory movements: %w", err)
	}
	defer rows.Close()

	return r.scanMovements(rows)
}

func (r *inventoryMovementRepository) scanMovements(rows *sql.Rows) ([]*entities.InventoryMovement, error) {
	var movements []*entities.InventoryMovement
	for rows.Next() {
		var m entities.InventoryMovement
		err := rows.Scan(
			&m.ID,
			&m.ProductID,
			&m.WarehouseID,
			&m.MovementType,
			&m.Quantity,
			&m.ReferenceType,
			&m.ReferenceID,
			&m.CreatedBy,
			&m.Notes,
			&m.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan inventory movement: %w", err)
		}
		movements = append(movements, &m)
	}
	return movements, nil
}

// Helper functions for nullable fields
func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullableInt(i int) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(i), Valid: true}
}
