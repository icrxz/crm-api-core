package application

import (
	"context"
	"errors"
	"strconv"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type caseService struct {
	customerService CustomerService
	userService     UserService
	caseRepository  domain.CaseRepository
	productService  ProductService
}

type CaseService interface {
	CreateCase(ctx context.Context, newCase domain.CreateCase) (string, error)
	GetCaseByID(ctx context.Context, caseID string) (*domain.Case, error)
	SearchCases(ctx context.Context, filters domain.CaseFilters) (domain.PagingResult[domain.Case], error)
	UpdateCase(ctx context.Context, caseID string, newCase domain.CaseUpdate) error
}

func NewCaseService(
	customerService CustomerService,
	caseRepository domain.CaseRepository,
	productService ProductService,
	userService UserService,
) CaseService {
	return &caseService{
		customerService: customerService,
		caseRepository:  caseRepository,
		productService:  productService,
		userService:     userService,
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

	user, err := c.userService.Search(ctx, domain.UserFilters{
		Region: []string{regionStringified},
		Role:   []string{string(domain.OPERATOR)},
	})
	if err != nil {
		var customErr *domain.CustomError
		if !errors.As(err, &customErr) || !customErr.IsNotFound() {
			return domain.NewValidationError("user not found", nil)
		}
	}

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
