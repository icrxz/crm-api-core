package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type caseRepository struct {
	client *sqlx.DB
}

func NewCaseRepository(client *sqlx.DB) domain.CaseRepository {
	return &caseRepository{
		client: client,
	}
}

func (r *caseRepository) Create(ctx context.Context, crmCase domain.Case) (string, error) {
	crmCaseDTO := mapCaseToCaseDTO(crmCase)

	_, err := r.client.NamedExecContext(
		ctx,
		"INSERT INTO cases "+
			"(case_id, contractor_id, customer_id, origin, type, subject, priority, status, due_date, created_by, created_at, updated_by, updated_at, external_reference, product_id, region, owner_id) "+
			"VALUES "+
			"(:case_id, :contractor_id, :customer_id, :origin, :type, :subject, :priority, :status, :due_date, :created_by, :created_at, :updated_by, :updated_at, :external_reference, :product_id, :region, :owner_id)",
		crmCaseDTO,
	)
	if err != nil {
		return "", err
	}

	return crmCase.CaseID, nil
}

func (r *caseRepository) GetByID(ctx context.Context, caseID string) (*domain.Case, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("caseID is required", nil)
	}

	var crmCaseDTO CaseDTO
	err := r.client.GetContext(ctx, &crmCaseDTO, "SELECT * FROM cases WHERE case_id=$1", caseID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no case found with this id", map[string]any{"case_id": caseID})
		}
		return nil, err
	}

	crmCase := mapCaseDTOToCase(crmCaseDTO)
	return &crmCase, nil
}

func (r *caseRepository) Search(ctx context.Context, filters domain.CaseFilters) ([]domain.Case, error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.ContractorID, whereQuery, whereArgs, "contractor_id")
	whereQuery, whereArgs = prepareInQuery(filters.OwnerID, whereQuery, whereArgs, "owner_id")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "customer_id")
	whereQuery, whereArgs = prepareInQuery(filters.PartnerID, whereQuery, whereArgs, "partner_id")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "status")
	whereQuery, whereArgs = prepareInQuery(filters.Region, whereQuery, whereArgs, "region")

	query := fmt.Sprintf("SELECT * FROM cases WHERE %s", strings.Join(whereQuery, " AND "))

	var foundCases []CaseDTO
	err := r.client.SelectContext(ctx, &foundCases, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	crmCases := mapCaseDTOsToCases(foundCases)

	return crmCases, nil
}

func (r *caseRepository) Update(ctx context.Context, crmCase domain.Case) error {
	crmCaseDTO := mapCaseToCaseDTO(crmCase)

	_, err := r.client.NamedExecContext(
		ctx,
		"UPDATE cases SET "+
			"contractor_id = :contractor_id, "+
			"customer_id = :customer_id, "+
			"partner_id = :partner_id, "+
			"owner_id = :owner_id, "+
			"origin = :origin, "+
			"type = :type, "+
			"subject = :subject, "+
			"priority = :priority, "+
			"status = :status, "+
			"due_date = :due_date, "+
			"updated_by = :updated_by, "+
			"updated_at = :updated_at, "+
			"closed_at = :closed_at, "+
			"target_date = :target_date "+
			"WHERE case_id = :case_id",
		crmCaseDTO,
	)
	if err != nil {
		return err
	}

	return nil
}
