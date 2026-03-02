-- Migration: Create purchases and purchase_items tables
-- Implements tasks #38 and #39

-- Purchases table (Stock IN transactions)
CREATE TABLE IF NOT EXISTS purchases (
    id            SERIAL PRIMARY KEY,
    warehouse_id  INT            NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    created_by    INT            NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Purchase items table (line items for each purchase)
CREATE TABLE IF NOT EXISTS purchase_items (
    id            SERIAL PRIMARY KEY,
    purchase_id   INT            NOT NULL REFERENCES purchases(id) ON DELETE CASCADE,
    product_id    INT            NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity      INT            NOT NULL CHECK (quantity > 0),
    unit_price    NUMERIC(10, 2) DEFAULT NULL,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_purchases_warehouse_id ON purchases(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_purchases_created_by ON purchases(created_by);
CREATE INDEX IF NOT EXISTS idx_purchases_created_at ON purchases(created_at);
CREATE INDEX IF NOT EXISTS idx_purchase_items_purchase_id ON purchase_items(purchase_id);
CREATE INDEX IF NOT EXISTS idx_purchase_items_product_id ON purchase_items(product_id);
