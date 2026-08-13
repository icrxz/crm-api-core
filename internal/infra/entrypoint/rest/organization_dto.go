package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateOrganizationDTO struct {
	Name      string `json:"name"`
	LegalName string `json:"legal_name"`
	Document  string `json:"document"`
	CreatedBy string `json:"created_by"`
}

type OrganizationDTO struct {
	OrganizationID string    `json:"organization_id"`
	Name           string    `json:"name"`
	LegalName      string    `json:"legal_name"`
	Document       string    `json:"document"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedBy      string    `json:"updated_by"`
	UpdatedAt      time.Time `json:"updated_at"`
	Active         bool      `json:"active"`
}

type UpdateOrganizationDTO struct {
	Name      *string `json:"name"`
	LegalName *string `json:"legal_name"`
	Document  *string `json:"document"`
	UpdatedBy string  `json:"updated_by"`
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

func mapCreateOrganizationDTOToOrganization(organizationDTO CreateOrganizationDTO) (domain.Organization, error) {
	return domain.NewOrganization(
		organizationDTO.Name,
		organizationDTO.LegalName,
		organizationDTO.Document,
		organizationDTO.CreatedBy,
	)
}

func mapOrganizationsToOrganizationDTOs(organizations []domain.Organization) []OrganizationDTO {
	organizationDTOs := make([]OrganizationDTO, 0, len(organizations))
	for _, organization := range organizations {
		organizationDTO := mapOrganizationToOrganizationDTO(organization)
		organizationDTOs = append(organizationDTOs, organizationDTO)
	}

	return organizationDTOs
}

func mapUpdateOrganizationDTOToUpdateOrganization(organizationDTO UpdateOrganizationDTO) domain.UpdateOrganization {
	return domain.UpdateOrganization{
		Name:      organizationDTO.Name,
		LegalName: organizationDTO.LegalName,
		Document:  organizationDTO.Document,
		UpdatedBy: organizationDTO.UpdatedBy,
	}
}
