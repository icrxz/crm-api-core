ALTER TABLE IF EXISTS partners
  ADD column IF NOT EXISTS active BOOLEAN DEFAULT true;

ALTER TABLE IF EXISTS contractors
  ADD column IF NOT EXISTS active BOOLEAN DEFAULT true;

ALTER TABLE IF EXISTS customers
  ADD column IF NOT EXISTS active BOOLEAN DEFAULT true;

CREATE INDEX IF NOT EXISTS idx_partners_active ON partners(active);
CREATE INDEX IF NOT EXISTS idx_contractors_active ON contractors(active);
CREATE INDEX IF NOT EXISTS idx_costumers_active ON customers(active);
