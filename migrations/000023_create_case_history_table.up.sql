CREATE TABLE IF NOT EXISTS case_history (
    history_id TEXT PRIMARY KEY,
    case_id TEXT NOT NULL REFERENCES cases(case_id),
    event_name VARCHAR(100) NOT NULL,
    author_id VARCHAR(100) NOT NULL,
    old_values JSONB NOT NULL DEFAULT '{}',
    new_values JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_case_history_case_id ON case_history (case_id);
