-- Rollback: Drop purchases and purchase_items tables

DROP INDEX IF EXISTS idx_purchase_items_product_id;
DROP INDEX IF EXISTS idx_purchase_items_purchase_id;
DROP INDEX IF EXISTS idx_purchases_created_at;
DROP INDEX IF EXISTS idx_purchases_created_by;
DROP INDEX IF EXISTS idx_purchases_warehouse_id;

DROP TABLE IF EXISTS purchase_items;
DROP TABLE IF EXISTS purchases;
