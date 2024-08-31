package database

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/icrxz/crm-api-core/internal/infra/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDatabase(config config.Database) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
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

	err = runMigrations(sqlClient)
	if err != nil {
		return nil, err
	}

	return sqlClient, nil
}

func runMigrations(client *sqlx.DB) error {
	driver, err := postgres.WithInstance(client.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	migration, err := migrate.NewWithDatabaseInstance("file:///app/migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return err
		}
	}

	return nil
}
