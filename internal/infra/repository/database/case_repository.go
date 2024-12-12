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

func (r *caseRepository) Search(ctx context.Context, filters domain.CaseFilters) (domain.PagingResult[domain.Case], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)
	limitArgs := make([]any, 0, 2)

	whereQuery, whereArgs = prepareInQuery(filters.ContractorID, whereQuery, whereArgs, "contractor_id")
	whereQuery, whereArgs = prepareInQuery(filters.OwnerID, whereQuery, whereArgs, "owner_id")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "customer_id")
	whereQuery, whereArgs = prepareInQuery(filters.PartnerID, whereQuery, whereArgs, "partner_id")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "status")
	whereQuery, whereArgs = prepareInQuery(filters.Region, whereQuery, whereArgs, "region")

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs = append(whereArgs, filters.Limit, filters.Offset)

	query := fmt.Sprintf("SELECT * FROM cases WHERE %s %s", strings.Join(whereQuery, " AND "), limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM cases WHERE %s", strings.Join(whereQuery, " AND "))

	var foundCases []CaseDTO
	err := r.client.SelectContext(ctx, &foundCases, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Case]{}, err
	}

	var countResult int
	err = r.client.GetContext(ctx, &countResult, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Case]{}, err
	}

	crmCases := mapCaseDTOsToCases(foundCases)

	result := domain.PagingResult[domain.Case]{
		Result: crmCases,
		Paging: domain.Paging{
			Total:  countResult,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
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

func (db *caseRepository) CreateBatch(ctx context.Context, cases []domain.Case) ([]string, error) {
	chunks := createChunks(cases, 100)
	tx := db.client.MustBegin()

	insertedIDs := make([]string, 0, len(cases))
	for _, chunk := range chunks {
		caseDTOs := mapCasesToCaseDTOs(chunk)

		query := "INSERT INTO cases " +
			"(case_id, contractor_id, customer_id, origin, type, subject, priority, status, due_date, created_by, created_at, updated_by, updated_at, external_reference, product_id, region, owner_id) " +
			"VALUES " +
			"(:case_id, :contractor_id, :customer_id, :origin, :type, :subject, :priority, :status, :due_date, :created_by, :created_at, :updated_by, :updated_at, :external_reference, :product_id, :region, :owner_id)" +
			"ON CONFLICT DO NOTHING"

		_, err := tx.NamedExecContext(
			ctx,
			query,
			caseDTOs,
		)
		if err != nil {
			return nil, err
		}

		for _, crmCase := range caseDTOs {
			insertedIDs = append(insertedIDs, crmCase.CaseID)
		}
	}

	err := tx.Commit()
	if err != nil {
		return nil, err
	}

	return insertedIDs, nil
}
