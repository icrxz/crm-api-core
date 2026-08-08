ALTER TABLE cases ADD COLUMN IF NOT EXISTS metadata JSONB NOT NULL DEFAULT '{}';

CREATE INDEX IF NOT EXISTS idx_cases_metadata ON cases USING GIN (metadata jsonb_path_ops);
