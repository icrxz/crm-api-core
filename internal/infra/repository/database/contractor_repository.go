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

type contractorRepository struct {
	client *sqlx.DB
}

func NewContractorRepository(client *sqlx.DB) domain.ContractorRepository {
	return &contractorRepository{
		client: client,
	}
}

func (db *contractorRepository) Create(ctx context.Context, contractor domain.Contractor) (string, error) {
	contractorDTO := mapContractorToContractorDTO(contractor)

	_, err := db.client.NamedExecContext(
		ctx,
		"INSERT INTO contractors "+
			"(contractor_id, company_name, legal_name, document, document_type, business_phone, business_email, created_at, created_by, updated_at, updated_by, active) "+
			"VALUES "+
			"(:contractor_id, :company_name, :legal_name, :document, :document_type, :business_phone, :business_email, :created_at, :created_by, :updated_at, :updated_by, :active)",
		contractorDTO,
	)
	if err != nil {
		return "", err
	}

	return contractor.ContractorID, nil
}

func (db *contractorRepository) Delete(ctx context.Context, contractorID string) error {
	if contractorID == "" {
		return domain.NewValidationError("contractorID is required", map[string]any{"contractor_id": contractorID})
	}

	_, err := db.client.ExecContext(ctx, "UPDATE contractors SET active = false WHERE contractor_id = $1", contractorID)
	if err != nil {
		return err
	}

	return nil
}

func (db *contractorRepository) GetByID(ctx context.Context, contractorID string) (*domain.Contractor, error) {
	var contractorDTO ContractorDTO
	err := db.client.GetContext(ctx, &contractorDTO, "SELECT * FROM contractors WHERE contractor_id=$1", contractorID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no contractor found with this id", map[string]any{"contractor_id": contractorID})
		}
		return nil, err
	}

	contractor := mapContractorDTOToContractor(contractorDTO)

	return &contractor, nil
}

func (db *contractorRepository) Search(ctx context.Context, filters domain.ContractorFilters) (domain.PagingResult[domain.Contractor], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)
	limitArgs := make([]any, 0, 2)

	whereQuery, whereArgs = prepareInQuery(filters.Document, whereQuery, whereArgs, "document")
	whereQuery, whereArgs = prepareInQuery(filters.CompanyName, whereQuery, whereArgs, "company_name")
	whereQuery, whereArgs = prepareInQuery(filters.ContractorID, whereQuery, whereArgs, "contractor_id")
	if filters.Active != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, strconv.FormatBool(*filters.Active))
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs = append(whereArgs, filters.Limit, filters.Offset)

	query := fmt.Sprintf("SELECT * FROM contractors WHERE %s %s", strings.Join(whereQuery, " AND "), limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM contractors WHERE %s", strings.Join(whereQuery, " AND "))

	var foundContractors []ContractorDTO
	err := db.client.SelectContext(ctx, &foundContractors, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.Contractor]{}, err
	}

	var countResult int
	err = db.client.GetContext(ctx, &countResult, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.Contractor]{}, err
	}

	contractors := mapContractorDTOsToContractors(foundContractors)

	result := domain.PagingResult[domain.Contractor]{
		Result: contractors,
		Paging: domain.Paging{
			Total:  countResult,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (db *contractorRepository) Update(ctx context.Context, contractor domain.Contractor) error {
	contractorDTO := mapContractorToContractorDTO(contractor)

	_, err := db.client.NamedExecContext(
		ctx,
		"UPDATE contractors SET "+
			"company_name = :company_name, "+
			"legal_name = :legal_name, "+
			"document = :document, "+
			"document_type = :document_type, "+
			"business_phone = :business_phone, "+
			"business_email = :business_email, "+
			"updated_at = :updated_at, "+
			"updated_by = :updated_by, "+
			"active = :active "+
			"WHERE contractor_id = :contractor_id",
		contractorDTO,
	)
	if err != nil {
		return err
	}

	return nil
}
