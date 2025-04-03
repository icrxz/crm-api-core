package application

import (
	"context"
	"errors"
	"strconv"

	"github.com/icrxz/crm-api-core/internal/domain"
	"golang.org/x/sync/errgroup"
)

type caseService struct {
	customerService    CustomerService
	userService        UserService
	caseRepository     domain.CaseRepository
	productService     ProductService
	commentService     CommentService
	transactionService TransactionService
	partnerService     PartnerService
	contractorService  ContractorService
}

type CaseService interface {
	CreateCase(ctx context.Context, newCase domain.CreateCase) (string, error)
	GetCaseByID(ctx context.Context, caseID string) (*domain.Case, error)
	SearchCases(ctx context.Context, filters domain.CaseFilters) (domain.PagingResult[domain.Case], error)
	UpdateCase(ctx context.Context, caseID string, newCase domain.CaseUpdate) error
	GetCaseFullByID(ctx context.Context, caseID string) (*domain.CaseFull, error)
}

func NewCaseService(
	customerService CustomerService,
	caseRepository domain.CaseRepository,
	productService ProductService,
	userService UserService,
	commentService CommentService,
	transactionService TransactionService,
	partnerService PartnerService,
	contractorService ContractorService,
) CaseService {
	return &caseService{
		customerService:    customerService,
		caseRepository:     caseRepository,
		productService:     productService,
		userService:        userService,
		commentService:     commentService,
		transactionService: transactionService,
		partnerService:     partnerService,
		contractorService:  contractorService,
	}
}

func (c *caseService) CreateCase(ctx context.Context, newCase domain.CreateCase) (string, error) {
	crmCase := newCase.Case
	customer, err := c.customerService.GetByID(ctx, crmCase.CustomerID)
	if err != nil {
		return "", err
	}
	crmCase.Region = customer.GetRegion()

	err = c.assignOwnerToNewCase(ctx, &crmCase)
	if err != nil {
		return "", err
	}

	productID, err := c.productService.CreateProduct(ctx, newCase.Product)
	if err != nil {
		return "", err
	}
	crmCase.ProductID = productID

	caseID, err := c.caseRepository.Create(ctx, crmCase)
	if err != nil {
		return "", err
	}

	return caseID, nil
}

func (c *caseService) GetCaseByID(ctx context.Context, caseID string) (*domain.Case, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("case id cannot be empty", nil)
	}

	return c.caseRepository.GetByID(ctx, caseID)
}

func (c *caseService) SearchCases(ctx context.Context, filters domain.CaseFilters) (domain.PagingResult[domain.Case], error) {
	return c.caseRepository.Search(ctx, filters)
}

func (c *caseService) assignOwnerToNewCase(ctx context.Context, crmCase *domain.Case) error {
	regionStringified := strconv.Itoa(crmCase.Region)

	searchResult, err := c.userService.Search(ctx, domain.UserFilters{
		Region: []string{regionStringified},
		Role:   []string{string(domain.OPERATOR), string(domain.ADMIN_OPERATOR)},
		PagingFilter: domain.PagingFilter{
			Limit:  1,
			Offset: 0,
		},
	})
	if err != nil {
		var customErr *domain.CustomError
		if !errors.As(err, &customErr) || !customErr.IsNotFound() {
			return domain.NewValidationError("user not found", nil)
		}
	}

	user := searchResult.Result

	if len(user) > 0 {
		crmCase.OwnerID = user[0].UserID
		crmCase.Status = domain.CUSTOMER_INFO
	}

	return nil
}

func (c *caseService) UpdateCase(ctx context.Context, caseID string, newCase domain.CaseUpdate) error {
	if caseID == "" {
		return domain.NewValidationError("case id cannot be empty", nil)
	}

	crmCase, err := c.caseRepository.GetByID(ctx, caseID)
	if err != nil {
		return err
	}

	crmCase.MergeUpdate(newCase)

	return c.caseRepository.Update(ctx, *crmCase)
}

func (c *caseService) GetCaseFullByID(ctx context.Context, caseID string) (*domain.CaseFull, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("case id cannot be empty", nil)
	}

	group := errgroup.Group{}

	var crmCase *domain.Case
	var comments []domain.Comment
	var transactions []domain.Transaction
	var product domain.Product
	var partner domain.Partner
	var customer domain.Customer
	var contractor domain.Contractor

	group.Go(func() error {
		foundCase, err := c.caseRepository.GetByID(ctx, caseID)
		if err != nil {
			return err
		}
		crmCase = foundCase
		return nil
	})

	group.Go(func() error {
		foundComments, err := c.commentService.GetByCaseID(ctx, caseID)
		if err != nil {
			return err
		}
		comments = foundComments
		return nil
	})

	group.Go(func() error {
		transactionFilter := domain.TransactionFilters{
			CaseIDs: []string{caseID},
		}

		foundTransactions, err := c.transactionService.SearchTransactions(ctx, transactionFilter)
		if err != nil {
			return err
		}

		transactions = foundTransactions

		return nil
	})

	groupErr := group.Wait()
	if groupErr != nil {
		return nil, groupErr
	}

	group.Go(func() error {
		foundProduct, err := c.productService.GetProductByID(ctx, crmCase.ProductID)
		if err != nil {
			var customErr *domain.CustomError
			if !errors.As(err, customErr) || !customErr.IsNotFound() {
				return err
			}
		}

		if foundProduct != nil {
			product = *foundProduct
		}

		return nil
	})

	group.Go(func() error {
		foundCustomer, err := c.customerService.GetByID(ctx, crmCase.CustomerID)
		if err != nil {
			var customErr *domain.CustomError
			if !errors.As(err, customErr) || !customErr.IsNotFound() {
				return err
			}
		}

		if foundCustomer != nil {
			customer = *foundCustomer
		}
		return nil
	})

	group.Go(func() error {
		foundPartner, err := c.partnerService.GetByID(ctx, crmCase.PartnerID)
		if err != nil {
			var customErr *domain.CustomError
			if !errors.As(err, customErr) || !customErr.IsNotFound() {
				return err
			}
		}

		if foundPartner != nil {
			partner = *foundPartner
		}

		return nil
	})

	group.Go(func() error {
		foundContractor, err := c.contractorService.GetByID(ctx, crmCase.ContractorID)
		if err != nil {
			var customErr *domain.CustomError
			if !errors.As(err, customErr) || !customErr.IsNotFound() {
				return err
			}
		}

		if foundContractor != nil {
			contractor = *foundContractor
		}

		return nil
	})

	groupErr = group.Wait()
	if groupErr != nil {
		return nil, groupErr
	}

	crmCaseFull := domain.NewCaseFull(*crmCase, comments, transactions, product, customer, partner, contractor)

	return &crmCaseFull, nil
}
