package application

import (
	"context"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type caseService struct {
	customerService CustomerService
	caseRepository  domain.CaseRepository
	productService  ProductService
}

type CaseService interface {
	CreateCase(ctx context.Context, newCase domain.CreateCase) (string, error)
	GetCaseByID(ctx context.Context, caseID string) (*domain.Case, error)
	SearchCases(ctx context.Context, filters domain.CaseFilters) ([]domain.Case, error)
}

func NewCaseService(customerService CustomerService, caseRepository domain.CaseRepository, productService ProductService) CaseService {
	return &caseService{
		customerService: customerService,
		caseRepository:  caseRepository,
		productService:  productService,
	}
}

func (c *caseService) CreateCase(ctx context.Context, newCase domain.CreateCase) (string, error) {
	crmCase := newCase.Case
	customer, err := c.customerService.GetByID(ctx, crmCase.CustomerID)
	if err != nil {
		return "", err
	}

	productID, err := c.productService.CreateProduct(ctx, newCase.Product)
	if err != nil {
		return "", err
	}

	crmCase.Region = customer.Region
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

func (c *caseService) SearchCases(ctx context.Context, filters domain.CaseFilters) ([]domain.Case, error) {
	return c.caseRepository.Search(ctx, filters)
}
