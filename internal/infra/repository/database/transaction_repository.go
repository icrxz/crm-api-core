package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type transactionRepository struct {
	client *sqlx.DB
}

func NewTransactionRepository(client *sqlx.DB) domain.TransactionRepository {
	return &transactionRepository{
		client: client,
	}
}

func (r *transactionRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) (string, error) {
	transactionDTO := mapTransactionToTransactionDTO(transaction)

	_, err := r.client.NamedExecContext(
		ctx,
		"INSERT INTO transactions "+
			"(transaction_id, case_id, type, amount, status, attachment_id, created_at, updated_at, created_by, updated_by, description) "+
			"VALUES "+
			"(:transaction_id, :case_id, :type, :amount, :status, :attachment_id, :created_at, :updated_at, :created_by, :updated_by, :description)",
		transactionDTO,
	)
	if err != nil {
		return "", err
	}

	return transaction.TransactionID, nil
}

func (r *transactionRepository) GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error) {
	if transactionID == "" {
		return domain.Transaction{}, domain.NewValidationError("transaction_id is required", nil)
	}

	var transactionDTO TransactionDTO
	err := r.client.GetContext(ctx, &transactionDTO, "SELECT * FROM transactions WHERE transaction_id=$1", transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Transaction{}, domain.NewNotFoundError("no transaction found with this id", map[string]any{"transaction_id": transactionID})
		}
		return domain.Transaction{}, err
	}

	return mapTransactionDTOToTransaction(transactionDTO), nil
}

func (r *transactionRepository) UpdateTransaction(ctx context.Context, transaction domain.Transaction) error {
	transactionDTO := mapTransactionToTransactionDTO(transaction)

	_, err := r.client.NamedExecContext(
		ctx,
		"UPDATE transactions SET "+
			"type = :type, "+
			"amount = :amount, "+
			"status = :status, "+
			"attachment_id = :attachment_id, "+
			"description = :description, "+
			"updated_at = :updated_at, "+
			"updated_by = :updated_by "+
			"WHERE transaction_id = :transaction_id",
		transactionDTO,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepository) SearchTransactions(ctx context.Context, filters domain.TransactionFilters) ([]domain.Transaction, error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.Types, whereQuery, whereArgs, "type")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "status")
	whereQuery, whereArgs = prepareInQuery(filters.CaseIDs, whereQuery, whereArgs, "case_id")

	query := fmt.Sprintf("SELECT * FROM transactions WHERE %s", strings.Join(whereQuery, " AND "))

	var foundTransactions []TransactionDTO
	err := r.client.SelectContext(ctx, &foundTransactions, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	transactions := mapTransactionDTOsToTransactions(foundTransactions)

	return transactions, nil
}

func (r *transactionRepository) CreateTransactionBatch(ctx context.Context, transactions []domain.Transaction) ([]string, error) {
	chunks := createChunks(transactions, 100)
	tx := r.client.MustBegin()

	insertedIDs := make([]string, 0, len(transactions))
	for _, chunk := range chunks {
		transactionDTOs := mapTransactionsToTransactionsDTOs(chunk)

		query := "INSERT INTO transactions " +
			"(transaction_id, case_id, type, amount, status, attachment_id, created_at, updated_at, created_by, updated_by, description) " +
			"VALUES " +
			"(:transaction_id, :case_id, :type, :amount, :status, :attachment_id, :created_at, :updated_at, :created_by, :updated_by, :description)"

		_, err := tx.NamedExecContext(
			ctx,
			query,
			transactionDTOs,
		)
		if err != nil {
			return nil, err
		}

		for _, transaction := range transactionDTOs {
			insertedIDs = append(insertedIDs, transaction.TransactionID)
		}
	}

	err := tx.Commit()
	if err != nil {
		return nil, err
	}

	return insertedIDs, nil
}

func (r *transactionRepository) DeleteManyByCaseID(ctx context.Context, caseID string) error {
	if caseID == "" {
		return domain.NewValidationError("caseID is required", nil)
	}

	_, err := r.client.ExecContext(ctx, "DELETE FROM transactions WHERE case_id = $1", caseID)
	if err != nil {
		return err
	}

	return nil
}
