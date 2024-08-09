package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type customerService struct {
	customerRepository domain.CustomerRepository
}

type CustomerService interface {
	Create(ctx context.Context, customer domain.Customer) (string, error)
	GetByID(ctx context.Context, customerID string) (*domain.Customer, error)
	Update(ctx context.Context, customerID string, updatedCustomer domain.UpdateCustomer) error
	Delete(ctx context.Context, customerID string) error
	Search(ctx context.Context, filters domain.CustomerFilters) (domain.PagingResult[domain.Customer], error)
}

func NewCustomerService(customerRepository domain.CustomerRepository) CustomerService {
	return &customerService{
		customerRepository: customerRepository,
	}
}

func (s *customerService) Create(ctx context.Context, customer domain.Customer) (string, error) {
	return s.customerRepository.Create(ctx, customer)
}

func (s *customerService) Delete(ctx context.Context, customerID string) error {
	if customerID == "" {
		return domain.NewValidationError("customerID cannot be empty", nil)
	}

	return s.customerRepository.Delete(ctx, customerID)
}

func (s *customerService) GetByID(ctx context.Context, customerID string) (*domain.Customer, error) {
	if customerID == "" {
		return nil, domain.NewValidationError("customerID cannot be empty", nil)
	}

	return s.customerRepository.GetByID(ctx, customerID)
}

func (s *customerService) Search(ctx context.Context, filters domain.CustomerFilters) (domain.PagingResult[domain.Customer], error) {
	return s.customerRepository.Search(ctx, filters)
}

func (s *customerService) Update(ctx context.Context, customerID string, updatedCustomer domain.UpdateCustomer) error {
	if customerID == "" {
		return domain.NewValidationError("customerID cannot be empty", nil)
	}

	customer, err := s.GetByID(ctx, customerID)
	if err != nil {
		return err
	}

	customer.MergeUpdate(updatedCustomer)

	return s.customerRepository.Update(ctx, *customer)
}
