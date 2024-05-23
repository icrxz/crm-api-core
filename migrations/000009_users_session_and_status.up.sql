ALTER TABLE IF EXISTS users
  ADD column IF NOT EXISTS active BOOLEAN DEFAULT true,
  ADD column IF NOT EXISTS last_logged_ip TEXT;

CREATE INDEX IF NOT EXISTS idx_users_active ON users(active);
