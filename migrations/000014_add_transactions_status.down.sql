ALTER TABLE IF EXISTS transactions
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS attachment_id,
    ALTER COLUMN type TYPE INTEGER USING (type::integer);

DROP INDEX IF EXISTS idx_transactions_status;
