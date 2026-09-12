package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type organizationRepository struct {
	client *sqlx.DB
}

func NewOrganizationRepository(client *sqlx.DB) domain.OrganizationRepository {
	return &organizationRepository{
		client: client,
	}
}

func (db *organizationRepository) Create(ctx context.Context, organization domain.Organization) (string, error) {
	organizationDTO := mapOrganizationToOrganizationDTO(organization)

	_, err := db.client.NamedExecContext(
		ctx,
		"INSERT INTO organizations "+
			"(organization_id, name, legal_name, document, created_at, created_by, updated_at, updated_by, active) "+
			"VALUES "+
			"(:organization_id, :name, :legal_name, :document, :created_at, :created_by, :updated_at, :updated_by, :active)",
		organizationDTO,
	)
	if err != nil {
		return "", err
	}

	return organization.OrganizationID, nil
}

func (db *organizationRepository) Delete(ctx context.Context, organizationID string) error {
	if organizationID == "" {
		return domain.NewValidationError("organizationID is required", map[string]any{"organization_id": organizationID})
	}

	_, err := db.client.ExecContext(ctx, "UPDATE organizations SET active = false WHERE organization_id = $1", organizationID)
	if err != nil {
		return err
	}

	return nil
}

func (db *organizationRepository) GetByID(ctx context.Context, organizationID string) (*domain.Organization, error) {
	var organizationDTO OrganizationDTO
	err := db.client.GetContext(ctx, &organizationDTO, "SELECT * FROM organizations WHERE organization_id=$1", organizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no organization found with this id", map[string]any{"organization_id": organizationID})
		}
		return nil, err
	}

	organization := mapOrganizationDTOToOrganization(organizationDTO)

	return &organization, nil
}

func (db *organizationRepository) Search(ctx context.Context, filters domain.OrganizationFilters) (domain.PagingResult[domain.Organization], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)
	var limitArgs []any

	whereQuery, whereArgs = prepareInQuery(filters.OrganizationID, whereQuery, whereArgs, "organization_id")
	whereQuery, whereArgs = prepareInQuery(filters.Name, whereQuery, whereArgs, "name")
	if filters.Active != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, strconv.FormatBool(*filters.Active))
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs = append(limitArgs, whereArgs...)
	limitArgs = append(limitArgs, filters.Limit, filters.Offset)

	query := fmt.Sprintf("SELECT * FROM organizations WHERE %s %s", strings.Join(whereQuery, " AND "), limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM organizations WHERE %s", strings.Join(whereQuery, " AND "))

	var foundOrganizations []OrganizationDTO
	err := db.client.SelectContext(ctx, &foundOrganizations, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Organization]{}, err
	}

	var countResult int
	err = db.client.GetContext(ctx, &countResult, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Organization]{}, err
	}

	organizations := mapOrganizationDTOsToOrganizations(foundOrganizations)

	result := domain.PagingResult[domain.Organization]{
		Result: organizations,
		Paging: domain.Paging{
			Total:  countResult,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (db *organizationRepository) Update(ctx context.Context, organization domain.Organization) error {
	organizationDTO := mapOrganizationToOrganizationDTO(organization)

	_, err := db.client.NamedExecContext(
		ctx,
		"UPDATE organizations SET "+
			"name = :name, "+
			"legal_name = :legal_name, "+
			"document = :document, "+
			"updated_at = :updated_at, "+
			"updated_by = :updated_by, "+
			"active = :active "+
			"WHERE organization_id = :organization_id",
		organizationDTO,
	)
	if err != nil {
		return err
	}

	return nil
}
