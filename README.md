# inventory_management_backend

## New: Sales (stock-out)

Endpoint:

- POST /sales

Payload example:

```json
{ "warehouse_id": 1, "product_id": 2, "quantity": 5, "unit_price": 9.99 }
```

Notes:

- This writes to `sales` and updates `product_stocks` atomically using a DB transaction.
- Create the following tables before using the endpoint:

```sql
CREATE TABLE product_stocks (
	product_id INT NOT NULL,
	warehouse_id INT NOT NULL,
	quantity INT NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(product_id, warehouse_id)
);

CREATE TABLE sales (
	id SERIAL PRIMARY KEY,
	warehouse_id INT NOT NULL,
	product_id INT NOT NULL,
	quantity INT NOT NULL,
	unit_price NUMERIC NOT NULL,
	created_at TIMESTAMP DEFAULT NOW()
);
```
