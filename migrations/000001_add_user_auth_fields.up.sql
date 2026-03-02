-- Migration: Add authentication and role-based access fields to users table

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS password_hash  TEXT        NOT NULL DEFAULT '',
    ADD COLUMN IF NOT EXISTS role           VARCHAR(20) NOT NULL DEFAULT 'staff',
    ADD COLUMN IF NOT EXISTS warehouse_id   INT         REFERENCES warehouses(id) ON DELETE SET NULL;

-- Set a check constraint to only allow valid roles
ALTER TABLE users
    ADD CONSTRAINT users_role_check CHECK (role IN ('admin', 'manager', 'staff'));

-- Add index on email for fast login lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Add index on role for fast RBAC filtering
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
