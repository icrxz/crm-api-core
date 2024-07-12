CREATE TABLE IF NOT EXISTS attachments (
    attachment_id TEXT primary key,
    comment_id TEXT,
    file_name TEXT,
    file_extension TEXT,
    attachment_url TEXT,
    created_at TIMESTAMP default now()
);