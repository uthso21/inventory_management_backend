-- Rollback: Drop inventory_movements table

DROP INDEX IF EXISTS idx_inventory_movements_reference;
DROP INDEX IF EXISTS idx_inventory_movements_created_at;
DROP INDEX IF EXISTS idx_inventory_movements_movement_type;
DROP INDEX IF EXISTS idx_inventory_movements_warehouse_id;
DROP INDEX IF EXISTS idx_inventory_movements_product_id;

DROP TABLE IF EXISTS inventory_movements;
