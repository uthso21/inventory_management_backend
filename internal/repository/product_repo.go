package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	entities "github.com/uthso21/inventory_management_backend/internal/entity"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductExists   = errors.New("product with this SKU already exists")
)

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	GetByID(ctx context.Context, id int) (*entities.Product, error)
	GetBySKU(ctx context.Context, sku string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entities.Product, error)
	GetSalesHistory(ctx context.Context, productID int) ([]*entities.SalesHistory, error)
	UpdateStock(ctx context.Context, productID int, quantity int) error
}

type productRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new instance of ProductRepository
func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *entities.Product) error {
	query := `
		INSERT INTO products (name, sku, description, category, unit_price, current_stock, reorder_point, lead_time_days, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	now := time.Now()
	product.CreatedAt = now
	product.UpdatedAt = now

	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.SKU,
		product.Description,
		product.Category,
		product.UnitPrice,
		product.CurrentStock,
		product.ReorderPoint,
		product.LeadTimeDays,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)

	if err != nil {
		if isDuplicateKeyError(err) {
			return ErrProductExists
		}
		return err
	}
	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id int) (*entities.Product, error) {
	query := `
		SELECT id, name, sku, description, category, unit_price, current_stock, reorder_point, lead_time_days, created_at, updated_at
		FROM products WHERE id = $1`

	product := &entities.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Description,
		&product.Category,
		&product.UnitPrice,
		&product.CurrentStock,
		&product.ReorderPoint,
		&product.LeadTimeDays,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *productRepository) GetBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	query := `
		SELECT id, name, sku, description, category, unit_price, current_stock, reorder_point, lead_time_days, created_at, updated_at
		FROM products WHERE sku = $1`

	product := &entities.Product{}
	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&product.ID,
		&product.Name,
		&product.SKU,
		&product.Description,
		&product.Category,
		&product.UnitPrice,
		&product.CurrentStock,
		&product.ReorderPoint,
		&product.LeadTimeDays,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product *entities.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, category = $3, unit_price = $4, current_stock = $5, reorder_point = $6, lead_time_days = $7, updated_at = $8
		WHERE id = $9`

	product.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Category,
		product.UnitPrice,
		product.CurrentStock,
		product.ReorderPoint,
		product.LeadTimeDays,
		product.UpdatedAt,
		product.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *productRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}

func (r *productRepository) List(ctx context.Context) ([]*entities.Product, error) {
	query := `
		SELECT id, name, sku, description, category, unit_price, current_stock, reorder_point, lead_time_days, created_at, updated_at
		FROM products ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*entities.Product
	for rows.Next() {
		product := &entities.Product{}
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.SKU,
			&product.Description,
			&product.Category,
			&product.UnitPrice,
			&product.CurrentStock,
			&product.ReorderPoint,
			&product.LeadTimeDays,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) GetSalesHistory(ctx context.Context, productID int) ([]*entities.SalesHistory, error) {
	query := `
		SELECT id, product_id, week_number, year, quantity_sold, revenue, created_at
		FROM sales_history
		WHERE product_id = $1
		ORDER BY year DESC, week_number DESC`

	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*entities.SalesHistory
	for rows.Next() {
		h := &entities.SalesHistory{}
		err := rows.Scan(
			&h.ID,
			&h.ProductID,
			&h.WeekNumber,
			&h.Year,
			&h.QuantitySold,
			&h.Revenue,
			&h.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return history, nil
}

func (r *productRepository) UpdateStock(ctx context.Context, productID int, quantity int) error {
	query := `UPDATE products SET current_stock = current_stock + $1, updated_at = $2 WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, quantity, time.Now(), productID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrProductNotFound
	}
	return nil
}
