CREATE TABLE IF NOT EXISTS customers (
    customer_id TEXT primary key,
    first_name TEXT,
    last_name TEXT,
    company_name TEXT,
    legal_name TEXT,
    customer_type TEXT,
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
    owner_id TEXT,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null
);

ALTER TABLE customers
    ADD CONSTRAINT fk_customers_owner_id
    FOREIGN KEY (owner_id)
    REFERENCES users(user_id);

CREATE INDEX IF NOT EXISTS idx_customers_owner_id ON customers(owner_id);
CREATE INDEX IF NOT EXISTS idx_customers_document ON customers(document);
CREATE INDEX IF NOT EXISTS idx_customers_customer_type ON customers(customer_type);
