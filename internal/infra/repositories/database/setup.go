package database

import "github.com/jmoiron/sqlx"

func NewDatabase() (*sqlx.DB, error) {
	sqlClient, err := sqlx.Connect("postgres", "user=postgres dbname=crm_users sslmode=disable")
	if err != nil {
		return nil, err
	}

	return sqlClient, nil
}
