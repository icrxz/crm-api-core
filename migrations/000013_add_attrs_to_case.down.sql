ALTER TABLE IF EXISTS cases
    DROP CONSTRAINT IF EXISTS fk_cases_product_id,
    DROP COLUMN IF EXISTS closed_at,
    DROP COLUMN IF EXISTS external_reference,
    DROP COLUMN IF EXISTS region,
    DROP COLUMN IF EXISTS product_id;
