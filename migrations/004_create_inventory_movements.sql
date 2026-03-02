-- Migration: Create inventory_movements table
-- Implements task #45 - Inventory movement log for stock tracking

CREATE TABLE IF NOT EXISTS inventory_movements (
    id              SERIAL PRIMARY KEY,
    product_id      INT            NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    warehouse_id    INT            NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    movement_type   VARCHAR(20)    NOT NULL CHECK (movement_type IN ('purchase', 'sale', 'adjustment', 'transfer')),
    quantity        INT            NOT NULL,
    reference_type  VARCHAR(50),   -- e.g., 'purchase', 'sale_order'
    reference_id    INT,           -- ID of the purchase or sale order
    created_by      INT            NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    notes           TEXT,
    created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

-- Indexes for inventory movement queries
CREATE INDEX IF NOT EXISTS idx_inventory_movements_product_id ON inventory_movements(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_warehouse_id ON inventory_movements(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_movement_type ON inventory_movements(movement_type);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_created_at ON inventory_movements(created_at);
CREATE INDEX IF NOT EXISTS idx_inventory_movements_reference ON inventory_movements(reference_type, reference_id);
