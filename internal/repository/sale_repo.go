package repository

import (
	"context"
	"database/sql"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

type SaleRepository interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, sale *entities.Sale) (int64, error)
}

type saleRepository struct{}

func NewSaleRepository() SaleRepository {
	return &saleRepository{}
}

func (r *saleRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, sale *entities.Sale) (int64, error) {

	query := `
	INSERT INTO sales (product_id, warehouse_id, quantity, sale_price, created_at)
	VALUES ($1,$2,$3,$4,NOW())
	RETURNING id
	`

	var id int64
	err := tx.QueryRowContext(ctx, query,
		sale.ProductID,
		sale.WarehouseID,
		sale.Quantity,
		sale.SalePrice,
	).Scan(&id)

	return id, err
}
