CREATE TABLE IF NOT EXISTS organizations (
    organization_id TEXT primary key,
    name TEXT,
    legal_name TEXT,
    document TEXT,
    created_at TIMESTAMP default now(),
    created_by TEXT not null,
    updated_at TIMESTAMP default now(),
    updated_by TEXT not null,
    active BOOLEAN default true
);

CREATE INDEX IF NOT EXISTS idx_organizations_name ON organizations(name);
