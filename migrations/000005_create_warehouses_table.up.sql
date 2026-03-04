-- Migration: Create warehouses table

CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    location TEXT,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
