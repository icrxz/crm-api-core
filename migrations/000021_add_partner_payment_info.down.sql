ALTER TABLE IF EXISTS partners
    DROP COLUMN IF EXISTS payment_type,
    DROP COLUMN IF EXISTS payment_owner,
    DROP COLUMN IF EXISTS payment_is_same_from_owner;
