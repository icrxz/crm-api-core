package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type organizationService struct {
	organizationRepository domain.OrganizationRepository
}

//go:generate mockgen -source=organization_service.go -destination=mock_application/mock_organization_service.go -package=mock_application
type OrganizationService interface {
	Create(ctx context.Context, organization domain.Organization) (string, error)
	GetByID(ctx context.Context, organizationID string) (*domain.Organization, error)
	Update(ctx context.Context, organizationID string, organization domain.UpdateOrganization) error
	Delete(ctx context.Context, organizationID string) error
	Search(ctx context.Context, filters domain.OrganizationFilters) (domain.PagingResult[domain.Organization], error)
}

func NewOrganizationService(organizationRepository domain.OrganizationRepository) OrganizationService {
	return &organizationService{
		organizationRepository: organizationRepository,
	}
}

func (s *organizationService) Create(ctx context.Context, organization domain.Organization) (string, error) {
	return s.organizationRepository.Create(ctx, organization)
}

func (s *organizationService) Delete(ctx context.Context, organizationID string) error {
	if organizationID == "" {
		return domain.NewValidationError("organizationID cannot be empty", nil)
	}

	return s.organizationRepository.Delete(ctx, organizationID)
}

func (s *organizationService) GetByID(ctx context.Context, organizationID string) (*domain.Organization, error) {
	if organizationID == "" {
		return nil, domain.NewValidationError("organizationID cannot be empty", nil)
	}

	return s.organizationRepository.GetByID(ctx, organizationID)
}

func (s *organizationService) Search(ctx context.Context, filters domain.OrganizationFilters) (domain.PagingResult[domain.Organization], error) {
	return s.organizationRepository.Search(ctx, filters)
}

func (s *organizationService) Update(ctx context.Context, organizationID string, updateOrganization domain.UpdateOrganization) error {
	if organizationID == "" {
		return domain.NewValidationError("organizationID cannot be empty", nil)
	}

	organization, err := s.GetByID(ctx, organizationID)
	if err != nil {
		return err
	}

	organization.MergeUpdate(updateOrganization)

	return s.organizationRepository.Update(ctx, *organization)
}
