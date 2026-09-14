package application

import (
	"context"
	"testing"

	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/icrxz/crm-api-core/internal/domain/mock_domain"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

const (
	testCaseID = "case-1"
	testAuthor = "user-1"
)

type caseActionServiceMocks struct {
	caseRepository        *mock_domain.MockCaseRepository
	caseHistoryRepository *mock_domain.MockCaseHistoryRepository
	transactionManager    *mock_domain.MockTransactionManager
	commentService        *mock_application.MockCommentService
	reportService         *mock_application.MockReportService
	transactionService    *mock_application.MockTransactionService
}

func newCaseActionServiceForTest(t *testing.T) (CaseActionService, *caseActionServiceMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mocks := &caseActionServiceMocks{
		caseRepository:        mock_domain.NewMockCaseRepository(ctrl),
		caseHistoryRepository: mock_domain.NewMockCaseHistoryRepository(ctrl),
		transactionManager:    mock_domain.NewMockTransactionManager(ctrl),
		commentService:        mock_application.NewMockCommentService(ctrl),
		reportService:         mock_application.NewMockReportService(ctrl),
		transactionService:    mock_application.NewMockTransactionService(ctrl),
	}

	mocks.transactionManager.EXPECT().
		WithinTransaction(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		AnyTimes()

	service := NewCaseActionService(
		mocks.caseRepository,
		mocks.caseHistoryRepository,
		mocks.transactionManager,
		mocks.commentService,
		mocks.reportService,
		nil,
		mocks.transactionService,
	)

	return service, mocks
}

func TestCaseActionService_ChangeStatus(t *testing.T) {
	t.Run("approves pending outgoing partner-payment transactions when the case closes", func(t *testing.T) {
		service, mocks := newCaseActionServiceForTest(t)

		crmCase := &domain.Case{CaseID: testCaseID, Status: domain.PAYMENT}
		mocks.caseRepository.EXPECT().GetByID(gomock.Any(), testCaseID).Return(crmCase, nil)
		mocks.caseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mocks.caseHistoryRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		transactions := []domain.Transaction{
			{TransactionID: "tx-mo", Status: domain.TRANSACTION_PENDING, Description: "MO"},
			{TransactionID: "tx-transport", Status: domain.TRANSACTION_PENDING, Description: "Deslocamento Técnico"},
			{TransactionID: "tx-parts", Status: domain.TRANSACTION_PENDING, Description: "Peças técnico"},
			{TransactionID: "tx-already-approved", Status: domain.TRANSACTION_APPROVED, Description: "MO"},
			{TransactionID: "tx-unrelated", Status: domain.TRANSACTION_PENDING, Description: "Cobrado seguradora"},
		}
		mocks.transactionService.EXPECT().
			SearchTransactions(gomock.Any(), domain.TransactionFilters{
				CaseIDs: []string{testCaseID},
				Types:   []string{string(domain.OUTGOING)},
			}).
			Return(transactions, nil)

		for _, txID := range []string{"tx-mo", "tx-transport", "tx-parts"} {
			mocks.transactionService.EXPECT().
				UpdateTransaction(gomock.Any(), txID, gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, update domain.TransactionUpdate) error {
					require.NotNil(t, update.Status)
					require.Equal(t, domain.TRANSACTION_APPROVED, *update.Status)
					return nil
				})
		}

		newStatus := domain.ChangeStatus{Status: domain.CLOSED, UpdatedBy: testAuthor}
		err := service.ChangeStatus(context.Background(), testCaseID, newStatus)

		require.NoError(t, err)
	})

	t.Run("does not touch transactions when the case does not close", func(t *testing.T) {
		service, mocks := newCaseActionServiceForTest(t)

		crmCase := &domain.Case{CaseID: testCaseID, Status: domain.PAYMENT}
		mocks.caseRepository.EXPECT().GetByID(gomock.Any(), testCaseID).Return(crmCase, nil)
		mocks.caseRepository.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		mocks.caseHistoryRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

		newStatus := domain.ChangeStatus{Status: domain.RECEIPT, UpdatedBy: testAuthor}
		err := service.ChangeStatus(context.Background(), testCaseID, newStatus)

		require.NoError(t, err)
	})
}
