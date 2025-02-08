package domain

type BuildCustomerFuncType func(row []string) (*Customer, error)

type BuildProductFuncType func(row []string) (*Product, error)

type CaseBuilder interface {
	GetCompanyName() []string
	GetCostumerDocumentIdx() int
	BuildCase(row []string, contractors []Contractor, customerID string, customerRegion int) (*Case, error)
	BuildProduct(row []string) (*Product, error)
	BuildCustomer(row []string) (*Customer, error)
}
