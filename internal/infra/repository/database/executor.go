package database

import (
	"context"
	"database/sql"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

// sqlExecutor is satisfied by both *sqlx.DB and *sqlx.Tx. Repositories should
// depend on this instead of calling their *sqlx.DB handle directly whenever
// they want to be usable from within a shared transaction.
type sqlExecutor interface {
	NamedExecContext(ctx context.Context, query string, arg any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...any) error
	SelectContext(ctx context.Context, dest any, query string, args ...any) error
	QueryxContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
}

type txKey struct{}

// executor returns the *sqlx.Tx bound to ctx by transactionManager.WithinTransaction,
// if any, falling back to db otherwise. Call this instead of using a repository's
// db handle directly for statements that must participate in the caller's transaction.
func executor(ctx context.Context, db *sqlx.DB) sqlExecutor {
	if tx, ok := ctx.Value(txKey{}).(*sqlx.Tx); ok {
		return tx
	}

	return db
}

type transactionManager struct {
	db *sqlx.DB
}

func NewTransactionManager(db *sqlx.DB) domain.TransactionManager {
	return &transactionManager{db: db}
}

func (t *transactionManager) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	if _, alreadyInTx := ctx.Value(txKey{}).(*sqlx.Tx); alreadyInTx {
		return fn(ctx)
	}

	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
