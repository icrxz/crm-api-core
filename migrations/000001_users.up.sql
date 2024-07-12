CREATE TABLE IF NOT EXISTS users (
  user_id TEXT primary key,
  first_name TEXT not null,
  last_name TEXT,
  username TEXT not null,
  email TEXT not null,
  password TEXT not null,
  role TEXT not null,
  region INT,
  created_at TIMESTAMP default now(),
  created_by TEXT not null,
  updated_at TIMESTAMP default now(),
  updated_by TEXT not null
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
