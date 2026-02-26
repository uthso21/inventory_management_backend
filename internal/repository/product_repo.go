package repository

import (
	"database/sql"

	"github.com/uthso21/inventory_management_backend/internal/entity"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (r *ProductRepository) CreateProduct(product *entity.Product) error {
	query := `INSERT INTO products (name, sku, price, description, stock) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.DB.Exec(query, product.Name, product.SKU, product.Price, product.Description, product.Stock)
	return err
}

func (r *ProductRepository) GetProducts() ([]entity.Product, error) {
	rows, err := r.DB.Query(`SELECT id, name, sku, price, description, stock FROM products`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []entity.Product

	for rows.Next() {
		var p entity.Product

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.SKU,
			&p.Price,
			&p.Description,
			&p.Stock,
		)

		if err = rows.Err(); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) DeleteProduct(id int) error {

	_, err := r.DB.Exec(`DELETE FROM products WHERE id = $1`, id)

	return err
}

func (r *ProductRepository) UpdateProduct(product *entity.Product) error {

	query := `
	UPDATE products
	SET name=$1, sku=$2, price=$3, description=$4, stock=$5
	WHERE id=$6
	`

	_, err := r.DB.Exec(
		query,
		product.Name,
		product.SKU,
		product.Price,
		product.Description,
		product.Stock,
		product.ID,
	)

	return err
}

func (r *ProductRepository) GetBySKU(sku string) (*entity.Product, error) {

	row := r.DB.QueryRow(
		`SELECT id, name, sku, price, description, stock FROM products WHERE sku=$1`,
		sku,
	)

	var p entity.Product

	err := row.Scan(
		&p.ID,
		&p.Name,
		&p.SKU,
		&p.Price,
		&p.Description,
		&p.Stock,
	)

	if err != nil {
		return nil, err
	}

	return &p, nil
}
