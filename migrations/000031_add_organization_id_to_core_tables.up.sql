ALTER TABLE users ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);
ALTER TABLE customers ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);
ALTER TABLE partners ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);
ALTER TABLE contractors ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);
ALTER TABLE cases ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);
ALTER TABLE products ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);
ALTER TABLE queues ADD COLUMN IF NOT EXISTS organization_id TEXT REFERENCES organizations(organization_id);

CREATE INDEX IF NOT EXISTS idx_users_organization_id ON users(organization_id);
CREATE INDEX IF NOT EXISTS idx_customers_organization_id ON customers(organization_id);
CREATE INDEX IF NOT EXISTS idx_partners_organization_id ON partners(organization_id);
CREATE INDEX IF NOT EXISTS idx_contractors_organization_id ON contractors(organization_id);
CREATE INDEX IF NOT EXISTS idx_cases_organization_id ON cases(organization_id);
CREATE INDEX IF NOT EXISTS idx_products_organization_id ON products(organization_id);
CREATE INDEX IF NOT EXISTS idx_queues_organization_id ON queues(organization_id);
