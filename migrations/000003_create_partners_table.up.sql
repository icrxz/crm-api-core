CREATE TABLE IF NOT EXISTS partners (
    partner_id TEXT primary key,
    first_name TEXT,
    last_name TEXT,
    company_name TEXT,
    legal_name TEXT,
    partner_type TEXT not null,
    document TEXT,
    document_type TEXT,
    shipping_address TEXT,
    shipping_city TEXT,
    shipping_state TEXT,
    shipping_country TEXT,
    shipping_zip_code TEXT,
    billing_address TEXT,
    billing_city TEXT,
    billing_state TEXT,
    billing_country TEXT,
    billing_zip_code TEXT,
    personal_phone TEXT,
    personal_email TEXT,
    business_phone TEXT,
    business_email TEXT,
    region INT,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);

CREATE INDEX IF NOT EXISTS idx_partners_region ON partners(region);
CREATE INDEX IF NOT EXISTS idx_partners_document ON partners(document);
CREATE INDEX IF NOT EXISTS idx_partners_partner_type ON partners(partner_type);
