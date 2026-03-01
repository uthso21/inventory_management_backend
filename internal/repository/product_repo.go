package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/uthso21/inventory_management_backend/internal/database"
	"github.com/uthso21/inventory_management_backend/internal/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductExists   = errors.New("product already exists")
)

type ProductRepository interface {
	Create(ctx context.Context, product *entity.Product) error
	GetByID(ctx context.Context, id int) (*entity.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entity.Product, error)
	Update(ctx context.Context, product *entity.Product) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entity.Product, error)
}

type productRepository struct {
	db *sql.DB
}

func NewProductRepository() ProductRepository {
	return &productRepository{db: database.DB}
}

func (r *productRepository) Create(ctx context.Context, product *entity.Product) error {
	query := `
		INSERT INTO products (name, sku, price, description, stock, reorder_level)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		product.Name,
		product.SKU,
		product.Price,
		product.Description,
		product.Stock,
		product.ReorderLevel,
	).Scan(&product.ID)
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*entity.Product, error) {
	query := `SELECT id, name, sku, price, description, stock, reorder_level FROM products WHERE id=$1`
	var p entity.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.SKU, &p.Price, &p.Description, &p.Stock, &p.ReorderLevel,
	)
	if err == sql.ErrNoRows {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entity.Product, error) {
	query := `SELECT id, name, sku, price, description, stock, reorder_level FROM products WHERE sku=$1`
	var p entity.Product
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&p.ID, &p.Name, &p.SKU, &p.Price, &p.Description, &p.Stock, &p.ReorderLevel,
	)
	if err == sql.ErrNoRows {
		return nil, ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *productRepository) Update(ctx context.Context, product *entity.Product) error {
	query := `
		UPDATE products
		SET name=$1, sku=$2, price=$3, description=$4, stock=$5, reorder_level=$6
		WHERE id=$7
	`
	res, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.SKU,
		product.Price,
		product.Description,
		product.Stock,
		product.ReorderLevel,
		product.ID,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *productRepository) List(ctx context.Context) ([]*entity.Product, error) {
	query := `SELECT id, name, sku, price, description, stock, reorder_level FROM products`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entity.Product
	for rows.Next() {
		var p entity.Product
		if err := rows.Scan(
			&p.ID, &p.Name, &p.SKU, &p.Price, &p.Description, &p.Stock, &p.ReorderLevel,
		); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}
	return products, nil
}
