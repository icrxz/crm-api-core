CREATE TABLE IF NOT EXISTS cases (
    case_id TEXT primary key,
    customer_id TEXT,
    contractor_id TEXT,
    partner_id TEXT,
    owner_id TEXT,
    origin TEXT,
    type TEXT,
    subject TEXT,
    status TEXT,
    priority TEXT,
    due_date TIMESTAMP,
    case_closed_at TIMESTAMP,
    case_closed_by TEXT,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);

ALTER TABLE cases
    ADD CONSTRAINT fk_cases_customer_id
        FOREIGN KEY (customer_id)
            REFERENCES customers(customer_id);

ALTER TABLE cases
    ADD CONSTRAINT fk_cases_contractor_id
        FOREIGN KEY (contractor_id)
            REFERENCES contractors(contractor_id);

ALTER TABLE cases
    ADD CONSTRAINT fk_cases_partner_id
        FOREIGN KEY (partner_id)
            REFERENCES partners(partner_id);

ALTER TABLE cases
    ADD CONSTRAINT fk_cases_owner_id
        FOREIGN KEY (owner_id)
            REFERENCES users(user_id);

CREATE INDEX IF NOT EXISTS idx_cases_status ON cases(status);
CREATE INDEX IF NOT EXISTS idx_cases_priority ON cases(priority);
CREATE INDEX IF NOT EXISTS idx_cases_due_date ON cases(due_date);
