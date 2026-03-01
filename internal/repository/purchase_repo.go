package repository

import (
	"context"
	"database/sql"
	"sync"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type PurchaseRepository interface {
	Create(ctx context.Context, purchase *entities.Purchase) (int64, error)
	CreateWithTx(ctx context.Context, tx *sql.Tx, purchase *entities.Purchase) (int64, error)
}

type purchaseRepository struct {
	mu        sync.Mutex
	purchases []entities.Purchase
}

func NewPurchaseRepository() PurchaseRepository {
	return &purchaseRepository{
		purchases: []entities.Purchase{},
	}
}

func (r *purchaseRepository) Create(ctx context.Context, purchase *entities.Purchase) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	purchase.ID = len(r.purchases) + 1
	r.purchases = append(r.purchases, *purchase)

	return int64(purchase.ID), nil
}

func (r *purchaseRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, purchase *entities.Purchase) (int64, error) {
	query := `
		INSERT INTO purchases (product_id, warehouse_id, quantity, purchase_price, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING id
	`

	var purchaseID int64
	err := tx.QueryRowContext(
		ctx,
		query,
		purchase.ProductID,
		purchase.WarehouseID,
		purchase.Quantity,
		purchase.PurchasePrice,
	).Scan(&purchaseID)
	if err != nil {
		return 0, err
	}

	purchase.ID = int(purchaseID)
	return purchaseID, nil
}
