package domain

import "context"

//go:generate mockgen -source=transaction_manager.go -destination=mock_domain/mock_transaction_manager.go -package=mock_domain
type TransactionManager interface {
	// WithinTransaction runs fn inside a single database transaction. Every
	// repository call made with the ctx passed to fn will participate in
	// that transaction, as long as the repository reads its executor from
	// ctx (see infra/repository/database.executor). Nested calls reuse the
	// outer transaction instead of opening a new one.
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}
