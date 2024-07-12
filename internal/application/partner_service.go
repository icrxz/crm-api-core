package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type partnerService struct {
	partnerRepository domain.PartnerRepository
}

type PartnerService interface {
	Create(ctx context.Context, partner domain.Partner) (string, error)
	GetByID(ctx context.Context, partnerID string) (*domain.Partner, error)
	Update(ctx context.Context, partnerID string, editPartner domain.EditPartner) error
	Delete(ctx context.Context, partnerID string) error
	Search(ctx context.Context, filters domain.PartnerFilters) ([]domain.Partner, error)
}

func NewPartnerService(partnerRepository domain.PartnerRepository) PartnerService {
	return &partnerService{
		partnerRepository: partnerRepository,
	}
}

func (s *partnerService) Create(ctx context.Context, partner domain.Partner) (string, error) {
	return s.partnerRepository.Create(ctx, partner)
}

func (s *partnerService) GetByID(ctx context.Context, partnerID string) (*domain.Partner, error) {
	if partnerID == "" {
		return nil, domain.NewValidationError("partnerID cannot be empty", nil)
	}

	return s.partnerRepository.GetByID(ctx, partnerID)
}

func (s *partnerService) Update(ctx context.Context, partnerID string, editPartner domain.EditPartner) error {
	if partnerID == "" {
		return domain.NewValidationError("partnerID cannot be empty", nil)
	}

	partner, err := s.partnerRepository.GetByID(ctx, partnerID)
	if err != nil {
		return err
	}

	partner.MergeUpdate(editPartner)

	return s.partnerRepository.Update(ctx, *partner)
}

func (s *partnerService) Delete(ctx context.Context, partnerID string) error {
	if partnerID == "" {
		return domain.NewValidationError("partnerID cannot be empty", nil)
	}

	return s.partnerRepository.Delete(ctx, partnerID)
}

func (s *partnerService) Search(ctx context.Context, filters domain.PartnerFilters) ([]domain.Partner, error) {
	return s.partnerRepository.Search(ctx, filters)
}
