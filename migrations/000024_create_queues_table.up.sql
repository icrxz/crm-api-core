CREATE TABLE IF NOT EXISTS queues (
    queue_id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    category TEXT NOT NULL,
    states TEXT[] NOT NULL DEFAULT '{}',
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    created_by TEXT NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_by TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_queues_states ON queues USING GIN (states);
CREATE INDEX IF NOT EXISTS idx_queues_category ON queues (category);
