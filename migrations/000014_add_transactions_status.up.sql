ALTER TABLE IF EXISTS transactions
    ADD COLUMN IF NOT EXISTS status TEXT NULL,
    ADD COLUMN IF NOT EXISTS attachment_id TEXT NULL,
    ALTER COLUMN type TYPE TEXT;

CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
