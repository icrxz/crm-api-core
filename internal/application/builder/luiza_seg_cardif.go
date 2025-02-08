package builder

import (
	"fmt"
	"slices"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type luizaSegBuilder struct {
	columnsIndex map[string]int
	author       string
	companyName  string
}

func NewLuizaSegBuilder(columnsIndex map[string]int, author, _ string) domain.CaseBuilder {
	return &luizaSegBuilder{
		columnsIndex: columnsIndex,
		author:       author,
	}
}

func (b *luizaSegBuilder) GetCompanyName() []string {
	return []string{"LuizaSeg", "Cardif"}
}

func (b *luizaSegBuilder) GetCostumerDocumentIdx() int {
	return -1
}

func (b *luizaSegBuilder) BuildCase(row []string, contractors []domain.Contractor, customerID string, customerRegion int) (*domain.Case, error) {
	dueDate := time.Now().Add(7 * 24 * time.Hour)

	contractorColumn := row[b.columnsIndex["BASE"]+1]
	var selectedContractorIdx int
	if contractorColumn == "Garantias" {
		selectedContractorIdx = slices.IndexFunc(contractors, func(c domain.Contractor) bool {
			return c.CompanyName == "Cardif"
		})
	} else {
		selectedContractorIdx = slices.IndexFunc(contractors, func(c domain.Contractor) bool {
			return c.CompanyName == "LuizaSeg"
		})
	}

	if selectedContractorIdx == -1 {
		return nil, domain.NewNotFoundError(fmt.Sprintf("Not found contractor for base %s", contractorColumn), nil)
	}

	selectedContractor := contractors[selectedContractorIdx]

	newCrmCase, err := domain.NewCase(
		selectedContractor.ContractorID,
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
