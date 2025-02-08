package builder

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type assurantBuilder struct {
	columnsIndex map[string]int
	author       string
	companyName  string
}

func NewAssurantBuilder(columnsIndex map[string]int, author, companyName string) domain.CaseBuilder {
	return &assurantBuilder{
		columnsIndex: columnsIndex,
		author:       author,
		companyName:  companyName,
	}
}

func (b *assurantBuilder) GetCompanyName() []string {
	return []string{b.companyName}
}

func (b *assurantBuilder) GetCostumerDocumentIdx() int {
	return b.columnsIndex["CPF Cliente"]
}

func (b *assurantBuilder) BuildCase(row []string, contractors []domain.Contractor, customerID string, customerRegion int) (*domain.Case, error) {
	dueDate := time.Now().Add(7 * 24 * time.Hour)

	newCrmCase, err := domain.NewCase(
		contractors[0].ContractorID,
		customerID,
		"csv",
		"insurance",
		row[b.columnsIndex["Defeito Reclamado"]],
		dueDate,
		b.author,
		row[b.columnsIndex["Número Sinistro"]],
	)
	if err != nil {
		return nil, err
	}

	newCrmCase.Region = customerRegion

	return &newCrmCase, nil
}

func (b *assurantBuilder) BuildProduct(row []string) (*domain.Product, error) {
	productValue := float64(0)
	var err error

	productValueStr := row[b.columnsIndex["Valor Produto"]]
	if productValueStr != "" {
		productValueParsed := strings.ReplaceAll(productValueStr, ",", "")
		productValue, err = strconv.ParseFloat(productValueParsed, 64)
		if err != nil {
			return nil, err
		}
	}

	newProduct, err := domain.NewProduct(
		"",
		"",
		productValue,
		row[b.columnsIndex["Marca"]],
		row[b.columnsIndex["Produto"]],
		row[b.columnsIndex["Número de Série"]],
		b.author,
	)
	if err != nil {
		return nil, err
	}

	return &newProduct, nil
}

func (b *assurantBuilder) BuildCustomer(row []string) (*domain.Customer, error) {
	clientName := row[b.columnsIndex["Nome Cliente"]]

	newCustomer, err := domain.NewCustomer(
		strings.Split(clientName, " ")[0],
		strings.Join(strings.Split(clientName, " ")[1:], " "),
		"",
		"",
		row[b.columnsIndex["CPF Cliente"]],
		string(domain.CPF),
		b.author,
		domain.Contact{
			PhoneNumber: row[b.columnsIndex["Telefone Celular"]],
			Email:       row[b.columnsIndex["E-mail"]],
		},
		domain.Contact{},
		domain.Address{
			Address: fmt.Sprintf("%s - %s", row[b.columnsIndex["Endereço"]], row[b.columnsIndex["Bairro"]]),
			City:    row[b.columnsIndex["Cidade"]],
			State:   domain.AcronymForState[row[b.columnsIndex["Estado"]]],
			Country: "brazil",
			ZipCode: row[b.columnsIndex["CEP"]],
		},
		domain.Address{},
	)
	if err != nil {
		return nil, err
	}

	return &newCustomer, nil
}
