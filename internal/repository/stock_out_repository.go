package repository

import (
    "database/sql"
    "errors"
    "github.com/uthso21/inventory_management_backend/internal/database"
)

type StockOutRepository struct {
    db *sql.DB
}

func NewStockOutRepository() *StockOutRepository {
    return &StockOutRepository{db: database.DB}
}

// Check current stock
func (r *StockOutRepository) GetInventory(productID, warehouseID int) (int, error) {
    var qty int
    err := r.db.QueryRow("SELECT quantity FROM inventory WHERE product_id=$1 AND warehouse_id=$2", productID, warehouseID).Scan(&qty)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, errors.New("no inventory found")
        }
        return 0, err
    }
    return qty, nil
}

// Reduce stock and insert stock_out transaction safely
func (r *StockOutRepository) StockOut(productID, warehouseID, quantity int, reason string) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    // Check stock
    var currentQty int
    err = tx.QueryRow("SELECT quantity FROM inventory WHERE product_id=$1 AND warehouse_id=$2 FOR UPDATE", productID, warehouseID).Scan(&currentQty)
    if err != nil {
        tx.Rollback()
        return err
    }

    if currentQty < quantity {
        tx.Rollback()
        return errors.New("not enough stock")
    }

    // Update inventory
    _, err = tx.Exec("UPDATE inventory SET quantity=quantity-$1, updated_at=NOW() WHERE product_id=$2 AND warehouse_id=$3", quantity, productID, warehouseID)
    if err != nil {
        tx.Rollback()
        return err
    }

    // Insert stock_out record
    _, err = tx.Exec("INSERT INTO stock_out (product_id, warehouse_id, quantity, reason, created_at) VALUES ($1,$2,$3,$4,NOW())", productID, warehouseID, quantity, reason)
    if err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}