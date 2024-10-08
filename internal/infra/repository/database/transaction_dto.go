package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type TransactionDTO struct {
	TransactionID string    `db:"transaction_id"`
	CaseID        string    `db:"case_id"`
	Type          string    `db:"type"`
	Value         float64   `db:"amount"`
	Status        string    `db:"status"`
	AttachmentID  string    `db:"attachment_id"`
	CreatedAt     time.Time `db:"created_at"`
	CreatedBy     string    `db:"created_by"`
	UpdatedAt     time.Time `db:"updated_at"`
	UpdatedBy     string    `db:"updated_by"`
	Description   *string   `db:"description"`
}

func mapTransactionToTransactionDTO(transaction domain.Transaction) TransactionDTO {
	return TransactionDTO{
		TransactionID: transaction.TransactionID,
		CaseID:        transaction.CaseID,
		Type:          string(transaction.Type),
		Value:         transaction.Value,
		Status:        string(transaction.Status),
		AttachmentID:  transaction.AttachmentID,
		CreatedAt:     transaction.CreatedAt,
		CreatedBy:     transaction.CreatedBy,
		UpdatedAt:     transaction.UpdatedAt,
		UpdatedBy:     transaction.UpdatedBy,
		Description:   &transaction.Description,
	}
}

func mapTransactionDTOToTransaction(transactionDTO TransactionDTO) domain.Transaction {
	var transactionDescription string
	if transactionDTO.Description != nil {
		transactionDescription = *transactionDTO.Description
	}

	return domain.Transaction{
		TransactionID: transactionDTO.TransactionID,
		CaseID:        transactionDTO.CaseID,
		Type:          domain.TransactionType(transactionDTO.Type),
		Value:         transactionDTO.Value,
		Status:        domain.TransactionStatus(transactionDTO.Status),
		AttachmentID:  transactionDTO.AttachmentID,
		CreatedAt:     transactionDTO.CreatedAt,
		CreatedBy:     transactionDTO.CreatedBy,
		UpdatedAt:     transactionDTO.UpdatedAt,
		UpdatedBy:     transactionDTO.UpdatedBy,
		Description:   transactionDescription,
	}
}

func mapTransactionDTOsToTransactions(transactionDTOs []TransactionDTO) []domain.Transaction {
	transactions := make([]domain.Transaction, 0, len(transactionDTOs))
	for _, transactionDTO := range transactionDTOs {
		transactions = append(transactions, mapTransactionDTOToTransaction(transactionDTO))
	}
	return transactions
}

func mapTransactionsToTransactionsDTOs(transactions []domain.Transaction) []TransactionDTO {
	transactionsDTOs := make([]TransactionDTO, 0, len(transactions))
	for _, transaction := range transactions {
		transactionsDTOs = append(transactionsDTOs, mapTransactionToTransactionDTO(transaction))
	}
	return transactionsDTOs
}
