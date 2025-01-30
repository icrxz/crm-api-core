package builder

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type luizaSegBuilder struct {
	columnsIndex map[string]int
	author       string
}

func NewLuizaSegBuilder(columnsIndex map[string]int, author string) domain.CaseBuilder {
	return &luizaSegBuilder{
		columnsIndex: columnsIndex,
		author:       author,
	}
}

func (b *luizaSegBuilder) GetCompanyName() string {
	return "LuizaSeg"
}

func (b *luizaSegBuilder) GetCostumerDocumentIdx() int {
	return -1
}

func (b *luizaSegBuilder) BuildCase(row []string, contractorID, customerID string, customerRegion int) (*domain.Case, error) {
	dueDate := time.Now().Add(7 * 24 * time.Hour)

	newCrmCase, err := domain.NewCase(
		contractorID,
		customerID,
		"csv",
		"insurance",
		"",
		dueDate,
		b.author,
		row[b.columnsIndex["SINISTRO"]+1],
	)
	if err != nil {
		return nil, err
	}

	newCrmCase.Region = customerRegion
	newCrmCase.Status = domain.DRAFT

	return &newCrmCase, nil
}

func (b *luizaSegBuilder) BuildProduct(row []string) (*domain.Product, error) {
	newProduct, err := domain.NewProduct(
		"",
		"",
		0.0,
		row[b.columnsIndex["MARCA"]+1],
		row[b.columnsIndex["PRODUTO"]+1],
		"",
		b.author,
	)
	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

func (b *luizaSegBuilder) BuildCustomer(row []string) (*domain.Customer, error) {
	newCustomer, err := domain.NewCustomer(
		"",
		"",
		"",
		"",
		"",
		string(domain.CPF),
		b.author,
		domain.Contact{},
		domain.Contact{},
		domain.Address{
			Address: "",
			City:    row[b.columnsIndex["CIDADE"]+1],
			State:   domain.AcronymForState[row[b.columnsIndex["UF"]+1]],
			Country: "brazil",
			ZipCode: "",
		},
		domain.Address{},
	)
	if err != nil {
		return nil, err
	}

	return &newCustomer, nil
}
