package database

import (
	"fmt"

	"github.com/icrxz/crm-api-core/internal/infra/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(config config.Database) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host(),
		config.Port,
		config.Username,
		config.Password(),
		config.Schema,
	)

	sqlClient, err := sqlx.Open(config.Driver, connectionString)
	if err != nil {
		return nil, err
	}

	sqlClient.SetMaxOpenConns(config.MaxOpenConns)
	sqlClient.SetMaxIdleConns(config.MaxIdleConns)
	sqlClient.SetConnMaxLifetime(config.ConnMaxLifetime)

	if err = sqlClient.Ping(); err != nil {
		sqlClient.Close()
		return nil, err
	}

	return sqlClient, nil
}
