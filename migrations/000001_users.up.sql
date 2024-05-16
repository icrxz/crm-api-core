CREATE TABLE IF NOT EXISTS users (
  user_id TEXT primary key,
  first_name TEXT not null,
  last_name TEXT,
  email TEXT not null,
  password TEXT not null,
  role TEXT not null,
  created_at TIMESTAMP default now(),
  created_by TEXT not null,
  updated_at TIMESTAMP default now(),
  updated_by TEXT not null
);