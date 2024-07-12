CREATE TABLE IF NOT EXISTS contractors_template (
    template_id TEXT PRIMARY KEY,
    contractor_id TEXT,
    template_url TEXT NOT NULL,
    template_login TEXT NOT NULL,
    template_password TEXT NOT NULL,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);

ALTER TABLE contractors_template ADD CONSTRAINT fk_contractor_id FOREIGN KEY (contractor_id) REFERENCES contractors(contractor_id);

CREATE INDEX IF NOT EXISTS idx_contractor_id ON contractors_template(contractor_id);
