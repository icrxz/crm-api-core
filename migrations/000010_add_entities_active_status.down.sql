DROP INDEX IF EXISTS idx_partners_active ON partners(active);
DROP INDEX IF EXISTS idx_contractors_active ON contractors(active);
DROP INDEX IF EXISTS idx_customers_active ON customers(active);

ALTER TABLE IF EXISTS partners
  DROP column IF EXISTS active;

ALTER TABLE IF EXISTS contractors
  DROP column IF EXISTS active;

ALTER TABLE IF EXISTS customers
  DROP column IF EXISTS active;
