package application

import (
	"context"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/icrxz/crm-api-core/internal/application/builder"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type batchCaseService struct {
	customerService   CustomerService
	productService    ProductService
	contractorService ContractorService
	caseRepository    domain.CaseRepository
}

type BatchCaseService interface {
	CreateBatch(ctx context.Context, file io.Reader, fileName, createdBy, company string) ([]string, error)
}

func NewBatchCaseService(customerService CustomerService, productService ProductService, contractorService ContractorService, caseRepository domain.CaseRepository) BatchCaseService {
	return &batchCaseService{
		customerService:   customerService,
		productService:    productService,
		contractorService: contractorService,
		caseRepository:    caseRepository,
	}
}

func (s *batchCaseService) CreateBatch(ctx context.Context, file io.Reader, fileName, createdBy, companyName string) ([]string, error) {
	fileNameSplit := strings.Split(fileName, ".")
	fileExtension := fileNameSplit[len(fileNameSplit)-1]

	var readFileFunc func(file io.Reader) ([][]string, error)
	if slices.Contains([]string{"csv"}, fileExtension) {
		readFileFunc = readCSV
	} else if slices.Contains([]string{"xls", "xlsx"}, fileExtension) {
		readFileFunc = readXLS
	}

	if readFileFunc == nil {
		return nil, domain.NewValidationError("file cannot be different from .csv, .xls", nil)
	}

	casesRows, err := readFileFunc(file)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err.Error())
		return nil, err
	}

	var caseBuilder domain.CaseBuilder
	var columnsIndex map[string]int
	var fileRows [][]string

	switch companyName {
	case "Assurant":
		columnsIndex = getColumnHeadersIndex(casesRows[0])
		fileRows = casesRows[1:]
		caseBuilder = builder.NewAssurantBuilder(columnsIndex, createdBy, companyName)
	case "Cardif":
		columnsIndex = getColumnHeadersIndex(casesRows[0])
		fileRows = casesRows[1:]
		caseBuilder = builder.NewLuizaSegBuilder(columnsIndex, createdBy, companyName)
	case "Ezze Seguros":
		columnsIndex = getColumnHeadersIndex(casesRows[0])
		fileRows = casesRows[1:]
		caseBuilder = builder.NewEzzeBuilder(columnsIndex, createdBy, companyName)
	case "LuizaSeg":
		columnsIndex = getColumnHeadersIndex(casesRows[0])
		fileRows = casesRows[1:]
		caseBuilder = builder.NewLuizaSegBuilder(columnsIndex, createdBy, companyName)
	default:
		columnsIndex = getColumnHeadersIndex(casesRows[0])
		fileRows = casesRows[1:]
		caseBuilder = builder.NewDefaultBuilder(columnsIndex, createdBy, companyName)
	}

	cases, err := s.buildCases(ctx, fileRows, caseBuilder)
	if err != nil {
		fmt.Printf("error building cases: %v\n", err.Error())
		return nil, err
	}

	caseIDs, err := s.caseRepository.CreateBatch(ctx, cases)
	if err != nil {
		fmt.Printf("error creating cases: %v\n", err.Error())
		return nil, err
	}

	return caseIDs, nil
}

func (s *batchCaseService) buildCases(ctx context.Context, csvRows [][]string, builder domain.CaseBuilder) ([]domain.Case, error) {
	crmCases := make([]domain.Case, 0, len(csvRows))

	contractors, err := s.getCompany(ctx, builder.GetCompanyName())
	if err != nil {
		fmt.Printf("error getting company: %v\n", err.Error())
		return nil, err
	}

	customers := make(map[string]*domain.Customer)
	if builder.GetCostumerDocumentIdx() != -1 {
		customers, err = s.getCustomers(ctx, csvRows, builder.GetCostumerDocumentIdx(), builder.BuildCustomer)
		if err != nil {
			fmt.Printf("error getting customers: %v\n", err.Error())
			return nil, err
		}
	}

	for _, row := range csvRows {
		if len(row) <= 1 {
			continue
		}

		customerID := ""
		customerRegion := -1
		customerDocIdx := builder.GetCostumerDocumentIdx()
		if customerDocIdx >= 0 {
			if len(row) > customerDocIdx {
				if customer, hasCustomer := customers[row[customerDocIdx]]; hasCustomer {
					customerID = customer.CustomerID
					customerRegion = customer.GetRegion()
				}
			}
		}

		newCrmCase, err := builder.BuildCase(row, contractors, customerID, customerRegion)
		if err != nil {
			fmt.Printf("error building case: %v\n", err.Error())
			return nil, err
		}

		productID, err := s.createProduct(ctx, row, builder.BuildProduct)
		if err != nil {
			fmt.Printf("error creating product: %v\n", err.Error())
			return nil, err
		}
		newCrmCase.ProductID = productID

		crmCases = append(crmCases, *newCrmCase)
	}

	return crmCases, nil
}

func (s *batchCaseService) searchCustomerBatch(ctx context.Context, customerDocument []string) (map[string]*domain.Customer, error) {
	filters := domain.CustomerFilters{Document: customerDocument, PagingFilter: domain.PagingFilter{Limit: 1000, Offset: 0}}

	customers, err := s.customerService.Search(ctx, filters)
	if err != nil {
		fmt.Printf("error searching customers: %v\n", err.Error())
		return nil, err
	}

	customersMap := make(map[string]*domain.Customer)
	for _, doc := range customerDocument {
		customerIdx := slices.IndexFunc(customers.Result, func(c domain.Customer) bool {
			return c.Document == doc
		})

		if customerIdx == -1 {
			customersMap[doc] = nil
		} else {
			customersMap[doc] = &customers.Result[customerIdx]
		}
	}

	return customersMap, nil
}

func (s *batchCaseService) getCompany(ctx context.Context, companyName []string) ([]domain.Contractor, error) {
	contractorFilters := domain.ContractorFilters{
		CompanyName: companyName,
		PagingFilter: domain.PagingFilter{
			Limit:  100,
			Offset: 0,
		},
	}

	contractorsResult, err := s.contractorService.Search(ctx, contractorFilters)
	if err != nil {
		fmt.Printf("error searching contractors: %v\n", err.Error())
		return nil, err
	}

	if contractorsResult.Paging.Total == 0 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("company not found with name %s", companyName), nil)
	}

	return contractorsResult.Result, nil
}

func (s *batchCaseService) getCustomers(ctx context.Context, rows [][]string, documentColumn int, buildCustomerFunc domain.BuildCustomerFuncType) (map[string]*domain.Customer, error) {
	customerDocuments := make([]string, 0)
	for _, row := range rows {
		customerDocument := row[documentColumn]

		if !slices.Contains(customerDocuments, customerDocument) {
			customerDocuments = append(customerDocuments, customerDocument)
		}
	}

	customers, err := s.searchCustomerBatch(ctx, customerDocuments)
	if err != nil {
		fmt.Printf("error searching customers: %v\n", err.Error())
		return nil, err
	}

	for customerDoc, customerVal := range customers {
		if customerVal != nil {
			continue
		}

		customerIdx := slices.IndexFunc(rows, func(row []string) bool {
			return row[documentColumn] == customerDoc
		})

		customerRow := rows[customerIdx]

		newCustomer, err := buildCustomerFunc(customerRow)
		if err != nil {
			fmt.Printf("error building customer: %v\n", err.Error())
			return nil, err
		}

		_, err = s.customerService.Create(ctx, *newCustomer)
		if err != nil {
			fmt.Printf("error creating customer: %v\n", err.Error())
			return nil, err
		}

		customers[newCustomer.Document] = newCustomer
	}

	return customers, nil
}

func (s *batchCaseService) createProduct(ctx context.Context, row []string, buildProductFunc domain.BuildProductFuncType) (string, error) {
	newProduct, err := buildProductFunc(row)
	if err != nil {
		fmt.Printf("error building product: %v\n", err.Error())
		return "", err
	}

	productID, err := s.productService.CreateProduct(ctx, *newProduct)
	if err != nil {
		fmt.Printf("error creating product: %v\n", err.Error())
		return "", err
	}

	return productID, nil
}
