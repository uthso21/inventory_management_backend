-- Insert test users
INSERT INTO users (username, email, password) VALUES
('admin', 'admin@inventory.com', '$2a$10$hashedpassword1'),
('john_doe', 'john@example.com', '$2a$10$hashedpassword2'),
('jane_smith', 'jane@example.com', '$2a$10$hashedpassword3')
ON CONFLICT (email) DO NOTHING;

-- Insert test products
INSERT INTO products (name, sku, description, category, unit_price, current_stock, reorder_point, lead_time_days) VALUES
('Laptop Pro 15', 'ELEC-001', 'High-performance laptop with 15-inch display', 'Electronics', 1299.99, 45, 20, 14),
('Wireless Mouse', 'ELEC-002', 'Ergonomic wireless mouse', 'Electronics', 29.99, 150, 50, 7),
('USB-C Hub', 'ELEC-003', '7-port USB-C hub with HDMI', 'Electronics', 49.99, 80, 30, 7),
('Office Chair', 'FURN-001', 'Ergonomic office chair with lumbar support', 'Furniture', 299.99, 25, 10, 21),
('Standing Desk', 'FURN-002', 'Electric height-adjustable desk', 'Furniture', 499.99, 15, 5, 28),
('Notebook Pack', 'STAT-001', 'Pack of 5 premium notebooks', 'Stationery', 12.99, 200, 100, 5),
('Pen Set', 'STAT-002', 'Professional pen set - 10 pcs', 'Stationery', 8.99, 300, 150, 5),
('Monitor 27"', 'ELEC-004', '4K monitor with HDR support', 'Electronics', 399.99, 35, 15, 14),
('Keyboard Mechanical', 'ELEC-005', 'RGB mechanical keyboard', 'Electronics', 89.99, 60, 25, 7),
('Webcam HD', 'ELEC-006', '1080p webcam with microphone', 'Electronics', 59.99, 100, 40, 7)
ON CONFLICT (sku) DO NOTHING;

-- Insert sales history for ML demand forecasting (last 12 weeks for product 1)
INSERT INTO sales_history (product_id, week_number, year, quantity_sold, revenue) VALUES
(1, 1, 2024, 12, 15599.88),
(1, 2, 2024, 8, 10399.92),
(1, 3, 2024, 15, 19499.85),
(1, 4, 2024, 10, 12999.90),
(1, 5, 2024, 18, 23399.82),
(1, 6, 2024, 14, 18199.86),
(1, 7, 2024, 20, 25999.80),
(1, 8, 2024, 16, 20799.84),
(1, 9, 2024, 22, 28599.78),
(1, 10, 2024, 19, 24699.81),
(1, 11, 2024, 25, 32499.75),
(1, 12, 2024, 21, 27299.79);

-- Insert sales history for product 2 (Wireless Mouse)
INSERT INTO sales_history (product_id, week_number, year, quantity_sold, revenue) VALUES
(2, 1, 2024, 45, 1349.55),
(2, 2, 2024, 52, 1559.48),
(2, 3, 2024, 38, 1139.62),
(2, 4, 2024, 60, 1799.40),
(2, 5, 2024, 55, 1649.45),
(2, 6, 2024, 48, 1439.52),
(2, 7, 2024, 70, 2099.30),
(2, 8, 2024, 65, 1949.35),
(2, 9, 2024, 58, 1739.42),
(2, 10, 2024, 72, 2159.28),
(2, 11, 2024, 80, 2399.20),
(2, 12, 2024, 75, 2249.25);

-- Insert some inventory transactions
INSERT INTO inventory_transactions (product_id, transaction_type, quantity, notes) VALUES
(1, 'IN', 50, 'Initial stock'),
(1, 'OUT', 5, 'Customer order #1001'),
(2, 'IN', 200, 'Initial stock'),
(2, 'OUT', 50, 'Customer order #1002'),
(3, 'IN', 100, 'Initial stock'),
(3, 'OUT', 20, 'Customer order #1003');
