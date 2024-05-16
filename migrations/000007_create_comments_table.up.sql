CREATE TABLE IF NOT EXISTS comments (
    comment_id TEXT primary key,
    case_id TEXT,
    content TEXT,
    comment_type TEXT,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);

ALTER TABLE comments
    ADD CONSTRAINT fk_comments_case_id
        FOREIGN KEY (case_id)
            REFERENCES cases(case_id);
