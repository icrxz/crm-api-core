package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type transactionService struct {
	transactionRepository domain.TransactionRepository
	caseService           CaseService
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction domain.Transaction) (string, error)
	GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error)
	UpdateTransaction(ctx context.Context, transactionID string, transactionUpdate domain.TransactionUpdate) error
	SearchTransactions(ctx context.Context, filters domain.TransactionFilters) ([]domain.Transaction, error)
	CreateTransactionBatch(ctx context.Context, transactions []domain.Transaction) ([]string, error)
}

func NewTransactionService(transactionRepository domain.TransactionRepository, caseService CaseService) TransactionService {
	return &transactionService{
		transactionRepository: transactionRepository,
		caseService:           caseService,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, transaction domain.Transaction) (string, error) {
	return s.transactionRepository.CreateTransaction(ctx, transaction)
}

func (s *transactionService) GetTransaction(ctx context.Context, transactionID string) (domain.Transaction, error) {
	if transactionID == "" {
		return domain.Transaction{}, domain.NewValidationError("transactionID cannot be empty", nil)
	}

	return s.transactionRepository.GetTransaction(ctx, transactionID)
}

func (s *transactionService) UpdateTransaction(ctx context.Context, transactionID string, transactionUpdate domain.TransactionUpdate) error {
	if transactionID == "" {
		return domain.NewValidationError("transactionID cannot be empty", nil)
	}

	transaction, err := s.transactionRepository.GetTransaction(ctx, transactionID)
	if err != nil {
		return err
	}

	transaction.MergeUpdate(transactionUpdate)

	return s.transactionRepository.UpdateTransaction(ctx, transaction)
}

func (s *transactionService) SearchTransactions(ctx context.Context, filters domain.TransactionFilters) ([]domain.Transaction, error) {
	return s.transactionRepository.SearchTransactions(ctx, filters)
}

func (s *transactionService) CreateTransactionBatch(ctx context.Context, transactions []domain.Transaction) ([]string, error) {
	if len(transactions) <= 0 || transactions[0].CaseID == "" {
		return nil, domain.NewValidationError("transactions cannot be empty and caseID cannot be empty", nil)
	}

	crmCase, err := s.caseService.GetCaseByID(ctx, transactions[0].CaseID)
	if err != nil {
		return nil, err
	}

	if crmCase.Status != domain.PAYMENT {
		return nil, domain.NewValidationError("case status must be in PAYMENT status", nil)
	}

	return s.transactionRepository.CreateTransactionBatch(ctx, transactions)
}
