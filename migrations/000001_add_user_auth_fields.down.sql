-- Rollback: Remove authentication and role-based access fields from users table

DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_users_email;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_role_check,
    DROP COLUMN IF EXISTS warehouse_id,
    DROP COLUMN IF EXISTS role,
    DROP COLUMN IF EXISTS password_hash;
