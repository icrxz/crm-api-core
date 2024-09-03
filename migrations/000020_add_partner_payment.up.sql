ALTER TABLE IF EXISTS partners
  ADD column IF NOT EXISTS payment_key TEXT,
  ADD column IF NOT EXISTS payment_key_option TEXT;
