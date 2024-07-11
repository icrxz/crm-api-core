package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction Transaction) (string, error)
	GetTransaction(ctx context.Context, transactionID string) (Transaction, error)
	UpdateTransaction(ctx context.Context, transaction Transaction) error
	SearchTransactions(ctx context.Context, filters TransactionFilters) ([]Transaction, error)
}

type Transaction struct {
	TransactionID string
	Type          TransactionType
	Value         float64
	CaseID        string
	AttachmentID  string
	Status        TransactionStatus
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedBy     string
	UpdatedAt     time.Time
	Description   string
}

type TransactionUpdate struct {
	Status       *TransactionStatus
	AttachmentID *string
	Value        *float64
	UpdatedBy    string
}

type TransactionFilters struct {
	CaseIDs []string
	Status  []string
	Types   []string
}

type TransactionType string

const (
	INCOMING TransactionType = "incoming"
	OUTGOING                 = "outgoing"
)

type TransactionStatus string

const (
	PENDING  TransactionStatus = "pending"
	APPROVED                   = "approved"
	REJECTED                   = "rejected"
)

func NewTransaction(
	transactionType TransactionType,
	value float64,
	caseID string,
	createdBy string,
	description string,
) (Transaction, error) {
	now := time.Now().UTC()

	transactionID, err := uuid.NewRandom()
	if err != nil {
		return Transaction{}, err
	}

	return Transaction{
		TransactionID: transactionID.String(),
		Type:          transactionType,
		Status:        PENDING,
		Value:         value,
		CaseID:        caseID,
		CreatedBy:     createdBy,
		CreatedAt:     now,
		UpdatedBy:     createdBy,
		UpdatedAt:     now,
		Description:   description,
	}, nil
}

func (t *Transaction) MergeUpdate(transactionUpdate TransactionUpdate) {
	t.UpdatedAt = time.Now().UTC()

	if transactionUpdate.Status != nil {
		t.Status = *transactionUpdate.Status
	}

	if transactionUpdate.AttachmentID != nil {
		t.AttachmentID = *transactionUpdate.AttachmentID
	}

	if transactionUpdate.Value != nil {
		t.Value = *transactionUpdate.Value
	}

	if transactionUpdate.UpdatedBy != "" {
		t.UpdatedBy = transactionUpdate.UpdatedBy
	}
}
