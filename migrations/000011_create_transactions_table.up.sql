CREATE TABLE IF NOT EXISTS transactions (
  transaction_id TEXT PRIMARY KEY,
  case_id TEXT NOT NULL,
  amount DECIMAL(10, 2) NOT NULL,
  type INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_by TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_by TEXT NOT NULL
);

ALTER TABLE transactions
    ADD CONSTRAINT fk_transactions_case_id
        FOREIGN KEY (case_id)
            REFERENCES cases(case_id);
