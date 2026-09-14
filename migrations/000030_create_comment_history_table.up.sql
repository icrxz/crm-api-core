CREATE TABLE IF NOT EXISTS comment_history (
    history_id TEXT PRIMARY KEY,
    comment_id TEXT NOT NULL REFERENCES comments(comment_id),
    event_name VARCHAR(100) NOT NULL,
    author_id VARCHAR(100) NOT NULL,
    old_values JSONB NOT NULL DEFAULT '{}',
    new_values JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP NOT NULL DEFAULT now()
);
