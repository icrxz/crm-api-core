package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type OrganizationDTO struct {
	OrganizationID string    `db:"organization_id"`
	Name           string    `db:"name"`
	LegalName      string    `db:"legal_name"`
	Document       string    `db:"document"`
	CreatedBy      string    `db:"created_by"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedBy      string    `db:"updated_by"`
	UpdatedAt      time.Time `db:"updated_at"`
	Active         bool      `db:"active"`
}

func mapOrganizationToOrganizationDTO(organization domain.Organization) OrganizationDTO {
	return OrganizationDTO{
		OrganizationID: organization.OrganizationID,
		Name:           organization.Name,
		LegalName:      organization.LegalName,
		Document:       organization.Document,
		CreatedBy:      organization.CreatedBy,
		CreatedAt:      organization.CreatedAt,
		UpdatedBy:      organization.UpdatedBy,
		UpdatedAt:      organization.UpdatedAt,
		Active:         organization.Active,
	}
}

func mapOrganizationDTOToOrganization(organizationDTO OrganizationDTO) domain.Organization {
	return domain.Organization{
		OrganizationID: organizationDTO.OrganizationID,
		Name:           organizationDTO.Name,
		LegalName:      organizationDTO.LegalName,
		Document:       organizationDTO.Document,
		CreatedBy:      organizationDTO.CreatedBy,
		CreatedAt:      organizationDTO.CreatedAt,
		UpdatedBy:      organizationDTO.UpdatedBy,
		UpdatedAt:      organizationDTO.UpdatedAt,
		Active:         organizationDTO.Active,
	}
}

func mapOrganizationDTOsToOrganizations(organizationDTOs []OrganizationDTO) []domain.Organization {
	organizations := make([]domain.Organization, 0, len(organizationDTOs))
	for _, organizationDTO := range organizationDTOs {
		organization := mapOrganizationDTOToOrganization(organizationDTO)
		organizations = append(organizations, organization)
	}

	return organizations
}
