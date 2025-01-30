package builder

import (
	"strconv"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type defaultBuilder struct {
	columnsIndex map[string]int
	author       string
	company      string
}

func NewDefaultBuilder(columnsIndex map[string]int, author string, company string) domain.CaseBuilder {
	return &defaultBuilder{
		columnsIndex: columnsIndex,
		author:       author,
		company:      company,
	}
}

func (b *defaultBuilder) GetCompanyName() string {
	return b.company
}

func (b *defaultBuilder) GetCostumerDocumentIdx() int {
	return b.columnsIndex["Documento"]
}

func (b *defaultBuilder) BuildCase(row []string, contractorID, customerID string, customerRegion int) (*domain.Case, error) {
	dueDate := time.Now().Add(7 * 24 * time.Hour)

	newCrmCase, err := domain.NewCase(
		contractorID,
		customerID,
		"csv",
		"insurance",
		row[b.columnsIndex["Descrição"]],
		dueDate,
		b.author,
		row[b.columnsIndex["Sinistro"]],
	)
	if err != nil {
		return nil, err
	}

	newCrmCase.Region = customerRegion

	return &newCrmCase, nil
}

func (b *defaultBuilder) BuildProduct(row []string) (*domain.Product, error) {
	productValue := float64(0)
	var err error
	productValueStr := row[b.columnsIndex["Valor"]]
	if productValueStr != "" {
		productValue, err = strconv.ParseFloat(productValueStr, 64)
		if err != nil {
			return nil, err
		}
	}

	newProduct, err := domain.NewProduct(
		"",
		"",
		productValue,
		row[b.columnsIndex["Marca"]],
		row[b.columnsIndex["Modelo"]],
		"",
		b.author,
	)
	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

func (b *defaultBuilder) BuildCustomer(row []string) (*domain.Customer, error) {
	newCustomer, err := domain.NewCustomer(
		row[b.columnsIndex["Nome"]],
		row[b.columnsIndex["Sobrenome"]],
		"",
		"",
		row[b.GetCostumerDocumentIdx()],
		string(domain.CPF),
		b.author,
		domain.Contact{},
		domain.Contact{},
		domain.Address{
			City:    row[b.columnsIndex["Cidade"]],
			State:   row[b.columnsIndex["Estado"]],
			Country: "brazil",
		},
		domain.Address{},
	)
	if err != nil {
		return nil, err
	}

	return &newCustomer, nil
}
