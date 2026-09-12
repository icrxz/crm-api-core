DROP INDEX IF EXISTS idx_users_organization_id;
DROP INDEX IF EXISTS idx_customers_organization_id;
DROP INDEX IF EXISTS idx_partners_organization_id;
DROP INDEX IF EXISTS idx_contractors_organization_id;
DROP INDEX IF EXISTS idx_cases_organization_id;
DROP INDEX IF EXISTS idx_products_organization_id;
DROP INDEX IF EXISTS idx_queues_organization_id;

ALTER TABLE users DROP COLUMN IF EXISTS organization_id;
ALTER TABLE customers DROP COLUMN IF EXISTS organization_id;
ALTER TABLE partners DROP COLUMN IF EXISTS organization_id;
ALTER TABLE contractors DROP COLUMN IF EXISTS organization_id;
ALTER TABLE cases DROP COLUMN IF EXISTS organization_id;
ALTER TABLE products DROP COLUMN IF EXISTS organization_id;
ALTER TABLE queues DROP COLUMN IF EXISTS organization_id;
