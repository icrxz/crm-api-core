package application

import (
	"context"
	"testing"

	"github.com/icrxz/crm-api-core/internal/application/mock_application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/icrxz/crm-api-core/internal/domain/mock_domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type caseServiceMocks struct {
	caseRepository        *mock_domain.MockCaseRepository
	caseHistoryRepository *mock_domain.MockCaseHistoryRepository
	customerService       *mock_application.MockCustomerService
	productService        *mock_application.MockProductService
	userService           *mock_application.MockUserService
	commentService        *mock_application.MockCommentService
	transactionService    *mock_application.MockTransactionService
	partnerService        *mock_application.MockPartnerService
	contractorService     *mock_application.MockContractorService
	queueService          *mock_application.MockQueueService
}

func newCaseServiceForTest(t *testing.T) (CaseService, *caseServiceMocks) {
	t.Helper()

	ctrl := gomock.NewController(t)

	mocks := &caseServiceMocks{
		caseRepository:        mock_domain.NewMockCaseRepository(ctrl),
		caseHistoryRepository: mock_domain.NewMockCaseHistoryRepository(ctrl),
		customerService:       mock_application.NewMockCustomerService(ctrl),
		productService:        mock_application.NewMockProductService(ctrl),
		userService:           mock_application.NewMockUserService(ctrl),
		commentService:        mock_application.NewMockCommentService(ctrl),
		transactionService:    mock_application.NewMockTransactionService(ctrl),
		partnerService:        mock_application.NewMockPartnerService(ctrl),
		contractorService:     mock_application.NewMockContractorService(ctrl),
		queueService:          mock_application.NewMockQueueService(ctrl),
	}

	service := NewCaseService(
		mocks.customerService,
		mocks.caseRepository,
		mocks.caseHistoryRepository,
		mock_domain.NewMockTransactionManager(ctrl),
		mocks.productService,
		mocks.userService,
		mocks.commentService,
		mocks.transactionService,
		mocks.partnerService,
		mocks.contractorService,
		mocks.queueService,
	)

	return service, mocks
}

func TestCaseService_GetCaseFullByID(t *testing.T) {
	t.Run("returns validation error when caseID is empty", func(t *testing.T) {
		service, _ := newCaseServiceForTest(t)

		_, err := service.GetCaseFullByID(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("resolves the full queue entity when the case has a queue_id", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)

		crmCase := &domain.Case{CaseID: "case-1", QueueID: "queue-1"}

		mocks.caseRepository.EXPECT().GetByID(gomock.Any(), "case-1").Return(crmCase, nil)
		mocks.commentService.EXPECT().GetByCaseID(gomock.Any(), "case-1").Return(nil, nil)
		mocks.transactionService.EXPECT().SearchTransactions(gomock.Any(), gomock.Any()).Return(nil, nil)
		mocks.queueService.EXPECT().GetByID(gomock.Any(), "queue-1").Return(&domain.Queue{QueueID: "queue-1", Name: "SP Mobile"}, nil)

		caseFull, err := service.GetCaseFullByID(context.Background(), "case-1")

		require.NoError(t, err)
		assert.Equal(t, "queue-1", caseFull.Queue.QueueID)
		assert.Equal(t, "SP Mobile", caseFull.Queue.Name)
	})

	t.Run("leaves queue empty when the case has no queue_id", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)

		crmCase := &domain.Case{CaseID: "case-1"}

		mocks.caseRepository.EXPECT().GetByID(gomock.Any(), "case-1").Return(crmCase, nil)
		mocks.commentService.EXPECT().GetByCaseID(gomock.Any(), "case-1").Return(nil, nil)
		mocks.transactionService.EXPECT().SearchTransactions(gomock.Any(), gomock.Any()).Return(nil, nil)

		caseFull, err := service.GetCaseFullByID(context.Background(), "case-1")

		require.NoError(t, err)
		assert.Equal(t, domain.Queue{}, caseFull.Queue)
	})

	t.Run("ignores a not-found queue instead of failing the whole request", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)

		crmCase := &domain.Case{CaseID: "case-1", QueueID: "queue-1"}

		mocks.caseRepository.EXPECT().GetByID(gomock.Any(), "case-1").Return(crmCase, nil)
		mocks.commentService.EXPECT().GetByCaseID(gomock.Any(), "case-1").Return(nil, nil)
		mocks.transactionService.EXPECT().SearchTransactions(gomock.Any(), gomock.Any()).Return(nil, nil)
		mocks.queueService.EXPECT().GetByID(gomock.Any(), "queue-1").
			Return(nil, domain.NewNotFoundError("no queue found with this id", nil))

		caseFull, err := service.GetCaseFullByID(context.Background(), "case-1")

		require.NoError(t, err)
		assert.Equal(t, domain.Queue{}, caseFull.Queue)
	})

	t.Run("propagates a non-not-found queue error", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)

		crmCase := &domain.Case{CaseID: "case-1", QueueID: "queue-1"}

		mocks.caseRepository.EXPECT().GetByID(gomock.Any(), "case-1").Return(crmCase, nil)
		mocks.commentService.EXPECT().GetByCaseID(gomock.Any(), "case-1").Return(nil, nil)
		mocks.transactionService.EXPECT().SearchTransactions(gomock.Any(), gomock.Any()).Return(nil, nil)
		mocks.queueService.EXPECT().GetByID(gomock.Any(), "queue-1").
			Return(nil, domain.NewValidationError("boom", nil))

		_, err := service.GetCaseFullByID(context.Background(), "case-1")

		require.Error(t, err)
	})
}

func TestCaseService_fetchRelatedEntities(t *testing.T) {
	t.Run("fetchProduct returns zero value on not-found", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)
		svc := service.(*caseService)

		mocks.productService.EXPECT().GetProductByID(gomock.Any(), "product-1").
			Return(nil, domain.NewNotFoundError("not found", nil))

		product, err := svc.fetchProduct(context.Background(), "product-1")

		require.NoError(t, err)
		assert.Equal(t, domain.Product{}, product)
	})

	t.Run("fetchProduct returns the found product", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)
		svc := service.(*caseService)

		mocks.productService.EXPECT().GetProductByID(gomock.Any(), "product-1").
			Return(&domain.Product{ProductID: "product-1"}, nil)

		product, err := svc.fetchProduct(context.Background(), "product-1")

		require.NoError(t, err)
		assert.Equal(t, "product-1", product.ProductID)
	})

	t.Run("fetchCustomer returns the found customer", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)
		svc := service.(*caseService)

		mocks.customerService.EXPECT().GetByID(gomock.Any(), "customer-1").
			Return(&domain.Customer{CustomerID: "customer-1"}, nil)

		customer, err := svc.fetchCustomer(context.Background(), "customer-1")

		require.NoError(t, err)
		assert.Equal(t, "customer-1", customer.CustomerID)
	})

	t.Run("fetchPartner returns the found partner", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)
		svc := service.(*caseService)

		mocks.partnerService.EXPECT().GetByID(gomock.Any(), "partner-1").
			Return(&domain.Partner{PartnerID: "partner-1"}, nil)

		partner, err := svc.fetchPartner(context.Background(), "partner-1")

		require.NoError(t, err)
		assert.Equal(t, "partner-1", partner.PartnerID)
	})

	t.Run("fetchContractor returns the found contractor", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)
		svc := service.(*caseService)

		mocks.contractorService.EXPECT().GetByID(gomock.Any(), "contractor-1").
			Return(&domain.Contractor{ContractorID: "contractor-1"}, nil)

		contractor, err := svc.fetchContractor(context.Background(), "contractor-1")

		require.NoError(t, err)
		assert.Equal(t, "contractor-1", contractor.ContractorID)
	})

	t.Run("propagates a non-not-found error from any related fetch", func(t *testing.T) {
		service, mocks := newCaseServiceForTest(t)
		svc := service.(*caseService)

		mocks.productService.EXPECT().GetProductByID(gomock.Any(), "product-1").
			Return(nil, domain.NewValidationError("boom", nil))

		_, err := svc.fetchProduct(context.Background(), "product-1")

		require.Error(t, err)
	})
}

func TestCaseService_SearchCasesFull(t *testing.T) {
	service, mocks := newCaseServiceForTest(t)

	filters := domain.CaseFilters{QueueID: []string{"queue-1"}}
	expected := domain.PagingResult[domain.CaseFull]{
		Result: []domain.CaseFull{{CaseID: "case-1"}},
	}

	mocks.caseRepository.EXPECT().SearchFull(gomock.Any(), filters).Return(expected, nil)

	result, err := service.SearchCasesFull(context.Background(), filters)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
