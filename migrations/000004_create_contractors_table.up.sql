CREATE TABLE IF NOT EXISTS contractors (
    contractor_id TEXT primary key,
    company_name TEXT,
    legal_name TEXT,
    document TEXT,
    document_type TEXT,
    business_phone TEXT,
    business_email TEXT,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);

CREATE INDEX IF NOT EXISTS idx_contractors_company_name ON contractors(company_name);
CREATE INDEX IF NOT EXISTS idx_contractors_document ON contractors(document);
