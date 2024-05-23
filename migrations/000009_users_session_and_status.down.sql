DROP INDEX IF EXISTS idx_users_active ON users(active);

ALTER TABLE IF EXISTS users
  DROP column IF EXISTS active,
  DROP column IF EXISTS last_logged_ip;
