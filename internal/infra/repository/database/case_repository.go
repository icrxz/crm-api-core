package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
	"strconv"
	"strings"
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

	result, err := r.client.NamedExecContext(
		ctx,
		"INSERT INTO cases "+
			"(case_id, contractor_id, customer_id, partner_id, owner_id, origin, type, subject, priority, status, due_date, created_by, created_at, updated_by, updated_at) "+
			"VALUES "+
			"(:case_id, :contractor_id, :customer_id, :partner_id, :owner_id, :origin, :type, :subject, :priority, :status, :due_date, :created_by, :created_at, :updated_by, :updated_at)",
		crmCaseDTO,
	)
	if err != nil {
		return "", err
	}

	caseID, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(caseID, 10), nil
}

func (r *caseRepository) GetByID(ctx context.Context, caseID string) (*domain.Case, error) {
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

func (r *caseRepository) Search(ctx context.Context) ([]domain.Case, error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

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
			"updated_at = :updated_at "+
			"WHERE case_id = :case_id",
		crmCaseDTO,
	)
	if err != nil {
		return err
	}

	return nil
}
