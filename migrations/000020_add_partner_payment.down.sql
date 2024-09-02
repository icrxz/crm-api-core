ALTER TABLE IF EXISTS partners
  DROP column IF EXISTS payment_key,
  DROP column IF EXISTS payment_key_option;