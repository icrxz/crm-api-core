ALTER TABLE IF EXISTS attachments
    ADD COLUMN created_by TEXT,
    ADD COLUMN key TEXT,
    ADD COLUMN size INTEGER;

CREATE INDEX IF NOT EXISTS idx_attachments_key ON attachments (key);