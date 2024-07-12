package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type contractorService struct {
	contractorRepository domain.ContractorRepository
}

type ContractorService interface {
	Create(ctx context.Context, contractor domain.Contractor) (string, error)
	GetByID(ctx context.Context, contractorID string) (*domain.Contractor, error)
	Update(ctx context.Context, contractorID string, contractor domain.UpdateContractor) error
	Delete(ctx context.Context, contractorID string) error
	Search(ctx context.Context, filters domain.ContractorFilters) ([]domain.Contractor, error)
}

func NewContractorService(contractorRepository domain.ContractorRepository) ContractorService {
	return &contractorService{
		contractorRepository: contractorRepository,
	}
}

func (s *contractorService) Create(ctx context.Context, contractor domain.Contractor) (string, error) {
	return s.contractorRepository.Create(ctx, contractor)
}

func (s *contractorService) Delete(ctx context.Context, contractorID string) error {
	if contractorID == "" {
		return domain.NewValidationError("contractorID cannot be empty", nil)
	}

	return s.contractorRepository.Delete(ctx, contractorID)
}

func (s *contractorService) GetByID(ctx context.Context, contractorID string) (*domain.Contractor, error) {
	if contractorID == "" {
		return nil, domain.NewValidationError("contractorID cannot be empty", nil)
	}

	return s.contractorRepository.GetByID(ctx, contractorID)
}

func (s *contractorService) Search(ctx context.Context, filters domain.ContractorFilters) ([]domain.Contractor, error) {
	return s.contractorRepository.Search(ctx, filters)
}

func (s *contractorService) Update(ctx context.Context, contractorID string, updateContractor domain.UpdateContractor) error {
	if contractorID == "" {
		return domain.NewValidationError("contractorID cannot be empty", nil)
	}

	contractor, err := s.GetByID(ctx, contractorID)
	if err != nil {
		return err
	}

	contractor.MergeUpdate(updateContractor)

	return s.contractorRepository.Update(ctx, *contractor)
}
