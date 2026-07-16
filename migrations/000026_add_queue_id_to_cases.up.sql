ALTER TABLE cases ADD COLUMN IF NOT EXISTS queue_id TEXT REFERENCES queues(queue_id);

CREATE INDEX IF NOT EXISTS idx_cases_queue_id ON cases (queue_id);
