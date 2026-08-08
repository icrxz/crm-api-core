DROP INDEX IF EXISTS idx_cases_metadata;
ALTER TABLE cases DROP COLUMN IF EXISTS metadata;
