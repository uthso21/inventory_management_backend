package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type MovementRepository interface {
	GetFiltered(ctx context.Context, mType string, productID string, date string) ([]*entities.Movement, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, movement *entities.Movement) error
}

type movementRepository struct{}

func NewMovementRepository() MovementRepository {
	return &movementRepository{}
}

func (r *movementRepository) GetFiltered(
	ctx context.Context,
	mType string,
	productID string,
	date string,
) ([]*entities.Movement, error) {

	query := `
	SELECT 
		m.id,
		m.product_id,
		m.warehouse_id,
		m.user_id,
		m.type,
		m.quantity,
		m.created_at,
		p.name as product_name,
		w.name as warehouse_name,
		COALESCE(u.name, '') as user_name
	FROM inventory_movements m
	JOIN products p ON m.product_id = p.id
	JOIN warehouses w ON m.warehouse_id = w.id
	LEFT JOIN users u ON m.user_id = u.id
	WHERE 1=1
	`

	args := []interface{}{}
	i := 1

	if mType != "" {
		query += " AND m.type=$" + fmt.Sprint(i)
		args = append(args, mType)
		i++
	}

	if productID != "" {
		query += " AND m.product_id=$" + fmt.Sprint(i)
		args = append(args, productID)
		i++
	}

	if date != "" {
		query += " AND DATE(m.created_at)=$" + fmt.Sprint(i)
		args = append(args, date)
		i++
	}

	query += " ORDER BY m.created_at DESC"

	rows, err := database.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movements []*entities.Movement

	for rows.Next() {
		var m entities.Movement
		err := rows.Scan(
			&m.ID,
			&m.ProductID,
			&m.WarehouseID,
			&m.UserID,
			&m.Type,
			&m.Quantity,
			&m.CreatedAt,
			&m.ProductName,
			&m.WarehouseName,
			&m.UserName,
		)
		if err != nil {
			return nil, err
		}

		movements = append(movements, &m)
	}

	return movements, nil
}

func (r *movementRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, movement *entities.Movement) error {

	query := `
	INSERT INTO inventory_movements 
	(product_id, warehouse_id, user_id, type, quantity, created_at)
	VALUES ($1, $2, $3, $4, $5, NOW())
	RETURNING id
	`

	return tx.QueryRowContext(
		ctx,
		query,
		movement.ProductID,
		movement.WarehouseID,
		movement.UserID,
		movement.Type,
		movement.Quantity,
	).Scan(&movement.ID)
}
