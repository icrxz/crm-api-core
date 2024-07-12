CREATE TABLE IF NOT EXISTS products (
    product_id TEXT primary key,
    name TEXT,
    description TEXT,
    brand TEXT,
    model TEXT,
    serial_number TEXT,
    value DECIMAL,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);
