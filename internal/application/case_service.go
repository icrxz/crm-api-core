package application

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"slices"
	"strconv"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type caseService struct {
	customerService   CustomerService
	userService       UserService
	caseRepository    domain.CaseRepository
	productService    ProductService
	contractorService ContractorService
}

type CaseService interface {
	CreateCase(ctx context.Context, newCase domain.CreateCase) (string, error)
	GetCaseByID(ctx context.Context, caseID string) (*domain.Case, error)
	SearchCases(ctx context.Context, filters domain.CaseFilters) (domain.PagingResult[domain.Case], error)
	UpdateCase(ctx context.Context, caseID string, newCase domain.CaseUpdate) error
	CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error)
}

func NewCaseService(
	customerService CustomerService,
	caseRepository domain.CaseRepository,
	productService ProductService,
	userService UserService,
	contractorService ContractorService,
) CaseService {
	return &caseService{
		customerService:   customerService,
		caseRepository:    caseRepository,
		productService:    productService,
		userService:       userService,
		contractorService: contractorService,
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
		Role:   []string{string(domain.OPERATOR)},
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

func (s *caseService) CreateBatch(ctx context.Context, file io.Reader, createdBy string) ([]string, error) {
	fileCSV := csv.NewReader(file)
	fileCSV.Comma = ';'

	casesRows, err := readCSV(fileCSV)
	if err != nil {
		return nil, err
	}

	columnsIndex := getColumnHeadersIndex(casesRows[0])
	cases, err := s.buildCases(ctx, casesRows[1:], columnsIndex, createdBy)
	if err != nil {
		return nil, err
	}

	return s.caseRepository.CreateBatch(ctx, cases)
}

func (s *caseService) buildCases(ctx context.Context, csvRows [][]string, columnsIndex map[string]int, author string) ([]domain.Case, error) {
	crmCases := make([]domain.Case, 0, len(csvRows))

	companyNames := make([]string, 0)
	customerDocuments := make([]string, 0)
	for _, row := range csvRows {
		companyName := row[columnsIndex["Seguradora"]]
		customerDocument := row[columnsIndex["Documento"]]

		if !slices.Contains(companyNames, companyName) {
			companyNames = append(companyNames, companyName)
		}

		if !slices.Contains(customerDocuments, customerDocument) {
			customerDocuments = append(customerDocuments, customerDocument)
		}
	}

	contractors, err := s.searchContractorBatch(ctx, companyNames)
	if err != nil {
		return nil, err
	}

	customers, err := s.searchCustomerBatch(ctx, customerDocuments)
	if err != nil {
		return nil, err
	}

	for _, row := range csvRows {
		customerDocument := row[columnsIndex["Documento"]]

		contractor := contractors[row[columnsIndex["Seguradora"]]]
		customer, hasCustomer := customers[customerDocument]
		if !hasCustomer {
			newCustomer, err := domain.NewCustomer(
				row[columnsIndex["Nome"]],
				row[columnsIndex["Sobrenome"]],
				"",
				"",
				customerDocument,
				string(domain.CPF),
				author,
				domain.Contact{},
				domain.Contact{},
				domain.Address{
					City:    row[columnsIndex["Cidade"]],
					State:   row[columnsIndex["Estado"]],
					Country: "brazil",
				},
				domain.Address{},
			)
			if err != nil {
				return nil, err
			}

			_, err = s.customerService.Create(ctx, newCustomer)
			if err != nil {
				return nil, err
			}

			customer = newCustomer
			customers[newCustomer.Document] = customer
		}

		dueDate := time.Now().Add(7 * 24 * time.Hour)

		newCrmCase, err := domain.NewCase(
			contractor.ContractorID,
			customer.CustomerID,
			"csv",
			"insurance",
			row[columnsIndex["Descrição"]],
			dueDate,
			author,
			row[columnsIndex["Sinistro"]],
		)
		if err != nil {
			return nil, err
		}

		newCrmCase.Region = customer.GetRegion()

		productValue := float64(0)
		productValueStr := row[columnsIndex["Valor"]]
		if productValueStr != "" {
			productValue, err = strconv.ParseFloat(productValueStr, 10)
			if err != nil {
				return nil, err
			}
		}

		newProduct, err := domain.NewProduct(
			"",
			"",
			productValue,
			row[columnsIndex["Marca"]],
			row[columnsIndex["Modelo"]],
			"",
			author,
		)
		if err != nil {
			return nil, err
		}

		productID, err := s.productService.CreateProduct(ctx, newProduct)
		if err != nil {
			return nil, err
		}
		newCrmCase.ProductID = productID

		crmCases = append(crmCases, newCrmCase)
	}
	return crmCases, nil
}

func (s *caseService) searchContractorBatch(ctx context.Context, companyName []string) (map[string]domain.Contractor, error) {
	filters := domain.ContractorFilters{CompanyName: companyName, PagingFilter: domain.PagingFilter{Limit: 1000, Offset: 0}}

	contractors, err := s.contractorService.Search(ctx, filters)
	if err != nil {
		return nil, err
	}

	if contractors.Paging.Total < len(companyName) {
		return nil, domain.NewNotFoundError(fmt.Sprintf("company not found with name %s", companyName), nil)
	}

	contractorsMap := make(map[string]domain.Contractor)
	for _, contractor := range contractors.Result {
		contractorsMap[contractor.CompanyName] = contractor
	}

	return contractorsMap, nil
}

func (s *caseService) searchCustomerBatch(ctx context.Context, customerDocument []string) (map[string]domain.Customer, error) {
	filters := domain.CustomerFilters{Document: customerDocument, PagingFilter: domain.PagingFilter{Limit: 1000, Offset: 0}}

	customers, err := s.customerService.Search(ctx, filters)
	if err != nil {
		return nil, err
	}

	customersMap := make(map[string]domain.Customer)
	for _, customer := range customers.Result {
		customersMap[customer.Document] = customer
	}

	return customersMap, nil
}
