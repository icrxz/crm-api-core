ALTER TABLE IF EXISTS cases
    ADD COLUMN IF NOT EXISTS closed_at TIMESTAMP NULL,
    ADD COLUMN IF NOT EXISTS external_reference TEXT NULL,
    ADD COLUMN IF NOT EXISTS region INT NULL,
    ADD COLUMN IF NOT EXISTS product_id TEXT NULL,
    ADD CONSTRAINT fk_cases_product_id
        FOREIGN KEY (product_id)
            REFERENCES products(product_id);
