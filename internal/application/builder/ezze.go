package builder

import (
	"strings"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type ezzeBuilder struct {
	columnsIndex map[string]int
	author       string
	companyName  string
}

func NewEzzeBuilder(columnsIndex map[string]int, author, companyName string) domain.CaseBuilder {
	return &ezzeBuilder{
		columnsIndex: columnsIndex,
		author:       author,
		companyName:  companyName,
	}
}

func (b *ezzeBuilder) GetCompanyName() []string {
	return []string{b.companyName}
}

func (b *ezzeBuilder) GetCostumerDocumentIdx() int {
	return b.columnsIndex["CPF/CNPJ Segurado"]
}

func (b *ezzeBuilder) BuildCase(row []string, contractors []domain.Contractor, customerID string, customerRegion int) (*domain.Case, error) {
	dueDate := time.Now().Add(7 * 24 * time.Hour)

	newCrmCase, err := domain.NewCase(
		contractors[0].ContractorID,
		customerID,
		"csv",
		"insurance",
		"",
		dueDate,
		b.author,
		row[b.columnsIndex["Ticket"]],
	)
	if err != nil {
		return nil, err
	}

	newCrmCase.Region = customerRegion
	newCrmCase.Status = domain.DRAFT

	return &newCrmCase, nil
}

func (b *ezzeBuilder) BuildProduct(row []string) (*domain.Product, error) {
	newProduct, err := domain.NewProduct(
		"",
		"",
		0.0,
		row[b.columnsIndex["Operação"]],
		row[b.columnsIndex["Bem Segurado"]],
		"",
		b.author,
	)
	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

func (b *ezzeBuilder) BuildCustomer(row []string) (*domain.Customer, error) {
	clientName := row[b.columnsIndex["Nome Segurado"]]

	newCustomer, err := domain.NewCustomer(
		strings.Split(clientName, " ")[0],
		strings.Join(strings.Split(clientName, " ")[1:], " "),
		"",
		"",
		row[b.GetCostumerDocumentIdx()],
		string(domain.CPF),
		b.author,
		domain.Contact{
			PhoneNumber: row[b.columnsIndex["Celular"]],
			Email:       row[b.columnsIndex["E-mail"]],
		},
		domain.Contact{},
		domain.Address{
			Address: "",
			City:    row[b.columnsIndex["Cidade"]],
			State:   domain.AcronymForState[row[b.columnsIndex["Estado"]]],
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
