package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type PurchaseRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, purchase *entities.Purchase) (int64, error)
	CreatePurchaseItemWithTx(ctx context.Context, tx *sql.Tx, item *entities.PurchaseItem) error
	GetByID(ctx context.Context, id int) (*entities.Purchase, error)
	List(ctx context.Context) ([]*entities.Purchase, error)
	GetItemsByPurchaseID(ctx context.Context, purchaseID int) ([]entities.PurchaseItem, error)
}

type purchaseRepository struct{}

func NewPurchaseRepository() PurchaseRepository {
	return &purchaseRepository{}
}

func (r *purchaseRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, purchase *entities.Purchase) (int64, error) {
	query := `
		INSERT INTO purchases (warehouse_id, created_by, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		purchase.WarehouseID,
		purchase.CreatedBy,
	).Scan(&purchase.ID, &purchase.CreatedAt, &purchase.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to create purchase: %w", err)
	}

	return int64(purchase.ID), nil
}

func (r *purchaseRepository) CreatePurchaseItemWithTx(ctx context.Context, tx *sql.Tx, item *entities.PurchaseItem) error {
	query := `
		INSERT INTO purchase_items (purchase_id, product_id, quantity, unit_price, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`

	err := tx.QueryRowContext(
		ctx,
		query,
		item.PurchaseID,
		item.ProductID,
		item.Quantity,
		item.UnitPrice,
	).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to create purchase item: %w", err)
	}

	return nil
}

func (r *purchaseRepository) GetByID(ctx context.Context, id int) (*entities.Purchase, error) {
	query := `
		SELECT id, warehouse_id, created_by, created_at, updated_at
		FROM purchases
		WHERE id = $1
	`

	var purchase entities.Purchase
	err := database.DB.QueryRowContext(ctx, query, id).Scan(
		&purchase.ID,
		&purchase.WarehouseID,
		&purchase.CreatedBy,
		&purchase.CreatedAt,
		&purchase.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("purchase not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase: %w", err)
	}

	// Get purchase items
	items, err := r.GetItemsByPurchaseID(ctx, id)
	if err != nil {
		return nil, err
	}
	purchase.Items = items

	return &purchase, nil
}

func (r *purchaseRepository) List(ctx context.Context) ([]*entities.Purchase, error) {
	query := `
		SELECT id, warehouse_id, created_by, created_at, updated_at
		FROM purchases
		ORDER BY created_at DESC
	`

	rows, err := database.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list purchases: %w", err)
	}
	defer rows.Close()

	var purchases []*entities.Purchase
	for rows.Next() {
		var purchase entities.Purchase
		err := rows.Scan(
			&purchase.ID,
			&purchase.WarehouseID,
			&purchase.CreatedBy,
			&purchase.CreatedAt,
			&purchase.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase: %w", err)
		}
		purchases = append(purchases, &purchase)
	}

	// Load items for each purchase
	for _, p := range purchases {
		items, err := r.GetItemsByPurchaseID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		p.Items = items
	}

	return purchases, nil
}

func (r *purchaseRepository) GetItemsByPurchaseID(ctx context.Context, purchaseID int) ([]entities.PurchaseItem, error) {
	query := `
		SELECT id, purchase_id, product_id, quantity, unit_price
		FROM purchase_items
		WHERE purchase_id = $1
		ORDER BY id
	`

	rows, err := database.DB.QueryContext(ctx, query, purchaseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get purchase items: %w", err)
	}
	defer rows.Close()

	var items []entities.PurchaseItem
	for rows.Next() {
		var item entities.PurchaseItem
		err := rows.Scan(
			&item.ID,
			&item.PurchaseID,
			&item.ProductID,
			&item.Quantity,
			&item.UnitPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan purchase item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
