package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/uthso21/inventory_management_backend/internal/database"
)

type StockOutRepository struct {
	db *sql.DB
}

func NewStockOutRepository() *StockOutRepository {
	return &StockOutRepository{db: database.DB}
}

// GetInventory returns current stock for product + warehouse
func (r *StockOutRepository) GetInventory(productID, warehouseID int) (int, error) {
	var qty int
	err := r.db.QueryRow("SELECT quantity FROM inventory WHERE product_id=$1 AND warehouse_id=$2",
		productID, warehouseID).Scan(&qty)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return qty, nil
}

// StockOut reduces stock and inserts a stock-out record safely
func (r *StockOutRepository) StockOut(productID, warehouseID, quantity int, reason string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	var currentQty int
	err = tx.QueryRow(
		"SELECT quantity FROM inventory WHERE product_id=$1 AND warehouse_id=$2 FOR UPDATE",
		productID, warehouseID,
	).Scan(&currentQty)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return errors.New("inventory not found")
		}
		return err
	}

	fmt.Printf("[DEBUG] ProductID=%d WarehouseID=%d CurrentQty=%d Requested=%d\n",
		productID, warehouseID, currentQty, quantity)

	if currentQty < quantity {
		tx.Rollback()
		return errors.New("not enough stock")
	}

	_, err = tx.Exec(
		"UPDATE inventory SET quantity = quantity - $1, updated_at = NOW() WHERE product_id=$2 AND warehouse_id=$3",
		quantity, productID, warehouseID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO stock_out (product_id, warehouse_id, quantity, reason, created_at) VALUES ($1,$2,$3,$4,NOW())",
		productID, warehouseID, quantity, reason,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}