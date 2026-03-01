CREATE TABLE IF NOT EXISTS products (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255)   NOT NULL,
    sku           VARCHAR(100)   NOT NULL UNIQUE,
    price         NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    description   TEXT,
    stock         INT            NOT NULL DEFAULT 0 CHECK (stock >= 0),
    reorder_level INT            NOT NULL DEFAULT 0 CHECK (reorder_level >= 0),
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ    NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);