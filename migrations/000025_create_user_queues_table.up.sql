CREATE TABLE IF NOT EXISTS user_queues (
    user_id TEXT NOT NULL REFERENCES users(user_id),
    queue_id TEXT NOT NULL REFERENCES queues(queue_id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, queue_id)
);

CREATE INDEX IF NOT EXISTS idx_user_queues_queue_id ON user_queues (queue_id);
