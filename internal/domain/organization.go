package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization Organization) (string, error)
	GetByID(ctx context.Context, organizationID string) (*Organization, error)
	Search(ctx context.Context, filters OrganizationFilters) (PagingResult[Organization], error)
	Update(ctx context.Context, organization Organization) error
	Delete(ctx context.Context, organizationID string) error
}

type Organization struct {
	OrganizationID string
	Name           string
	LegalName      string
	Document       string
	CreatedBy      string
	CreatedAt      time.Time
	UpdatedBy      string
	UpdatedAt      time.Time
	Active         bool
}

type UpdateOrganization struct {
	Name      *string
	LegalName *string
	Document  *string
	UpdatedBy string
}

func (o *Organization) MergeUpdate(newOrganization UpdateOrganization) {
	now := time.Now().UTC()
	o.UpdatedAt = now
	o.UpdatedBy = newOrganization.UpdatedBy

	if newOrganization.Name != nil {
		o.Name = *newOrganization.Name
	}

	if newOrganization.LegalName != nil {
		o.LegalName = *newOrganization.LegalName
	}

	if newOrganization.Document != nil {
		o.Document = *newOrganization.Document
	}
}

type OrganizationFilters struct {
	OrganizationID []string
	Name           []string
	Active         *bool
	PagingFilter
}

func NewOrganization(name, legalName, document, author string) (Organization, error) {
	now := time.Now().UTC()

	organizationID, err := uuid.NewRandom()
	if err != nil {
		return Organization{}, err
	}

	return Organization{
		OrganizationID: organizationID.String(),
		Name:           name,
		LegalName:      legalName,
		Document:       document,
		CreatedBy:      author,
		CreatedAt:      now,
		UpdatedBy:      author,
		UpdatedAt:      now,
		Active:         true,
	}, nil
}
