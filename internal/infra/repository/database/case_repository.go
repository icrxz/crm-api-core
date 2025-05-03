package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
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
	var limitArgs []any

	whereQuery, whereArgs = prepareInQuery(filters.ContractorID, whereQuery, whereArgs, "contractor_id")
	whereQuery, whereArgs = prepareInQuery(filters.OwnerID, whereQuery, whereArgs, "owner_id")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "customer_id")
	whereQuery, whereArgs = prepareInQuery(filters.PartnerID, whereQuery, whereArgs, "partner_id")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "status")
	whereQuery, whereArgs = prepareInQuery(filters.Region, whereQuery, whereArgs, "region")
	whereQuery, whereArgs = prepareLikeQuery(filters.ExternalReference, whereQuery, whereArgs, "external_reference")

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs = append(whereArgs, filters.Limit, filters.Offset)

	query := fmt.Sprintf("SELECT * FROM cases WHERE %s ORDER BY created_at DESC %s", strings.Join(whereQuery, " AND "), limitQuery)
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

func (r *caseRepository) CreateBatch(ctx context.Context, cases []domain.Case) ([]string, error) {
	chunks := createChunks(cases, 100)
	tx := r.client.MustBegin()

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

func (r *caseRepository) SearchFull(ctx context.Context, filters domain.CaseFilters) (domain.PagingResult[domain.CaseFull], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)
	var limitArgs []any

	whereQuery, whereArgs = prepareInQuery(filters.ContractorID, whereQuery, whereArgs, "ca.contractor_id")
	whereQuery, whereArgs = prepareInQuery(filters.OwnerID, whereQuery, whereArgs, "ca.owner_id")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "ca.customer_id")
	whereQuery, whereArgs = prepareInQuery(filters.PartnerID, whereQuery, whereArgs, "ca.partner_id")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "ca.status")
	whereQuery, whereArgs = prepareInQuery(filters.Region, whereQuery, whereArgs, "ca.region")
	whereQuery, whereArgs = prepareLikeQuery(filters.ExternalReference, whereQuery, whereArgs, "ca.external_reference")
	whereQuery, whereArgs = prepareInQuery(filters.ShippingState, whereQuery, whereArgs, "cu.shipping_state")

	if filters.StartDate != nil {
		whereQuery, whereArgs = prepareLesserEqualQuery(filters.StartDate, whereQuery, whereArgs, "ca.created_at")
	}

	if filters.EndDate != nil {
		whereQuery, whereArgs = prepareGreaterEqualQuery(filters.EndDate, whereQuery, whereArgs, "ca.created_at")
	}

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs = append(whereArgs, filters.Limit, filters.Offset)

	query := fmt.Sprintf(`SELECT
		ca.case_id,
		ca.owner_id,
		ca.origin,
		ca.type,
		ca.subject,
		ca.priority,
		ca.status,
		ca.due_date,
		ca.created_by,
		ca.created_at,
		ca.updated_by,
		ca.updated_at,
		ca.closed_at,
		ca.target_date,
		ca.external_reference,
		ca.region,
		co.contractor_id,
		co.company_name,
		co.legal_name,
		co.document,
		co.document_type,
		co.business_phone,
		co.business_email,
		co.created_by,
		co.created_at,
		co.updated_by,
		co.updated_at,
		co.active,
		cu.customer_id,
		cu.first_name,
		cu.last_name,
		cu.company_name,
		cu.legal_name,
		cu.customer_type,
		cu.document,
		cu.document_type,
		cu.shipping_address,
		cu.shipping_city,
		cu.shipping_state,
		cu.shipping_zip_code,
		cu.shipping_country,
		cu.billing_address,
		cu.billing_city,
		cu.billing_state,
		cu.billing_zip_code,
		cu.billing_country,
		cu.personal_phone,
		cu.business_phone,
		cu.personal_email,
		cu.business_email,
		cu.owner_id,
		cu.created_by,
		cu.created_at,
		cu.updated_by,
		cu.updated_at,
		cu.active,
		pa.partner_id,
		pa.first_name,
		pa.last_name,
		pa.company_name,
		pa.legal_name,
		pa.partner_type,
		pa.document,
		pa.document_type,
		pa.shipping_address,
		pa.shipping_city,
		pa.shipping_state,
		pa.shipping_zip_code,
		pa.shipping_country,
		pa.billing_address,
		pa.billing_city,
		pa.billing_state,
		pa.billing_zip_code,
		pa.billing_country,
		pa.personal_phone,
		pa.business_phone,
		pa.personal_email,
		pa.business_email,
		pa.region,
		pa.created_by,
		pa.created_at,
		pa.updated_by,
		pa.updated_at,
		pa.active,
		pa.description,
		pa.payment_key,
		pa.payment_key_option,
		pa.payment_type,
		pa.payment_owner,
		pa.payment_is_same_from_owner,
		pr.product_id,
		pr.name,
		pr.description,
		pr.brand,
		pr.model,
		pr.value,
		pr.serial_number,
		pr.created_at,
		pr.updated_at,
		pr.created_by,
		pr.updated_by
		FROM cases AS ca
		LEFT JOIN contractors AS co ON ca.contractor_id = co.contractor_id
		LEFT JOIN partners AS pa ON ca.partner_id = pa.partner_id
		LEFT JOIN customers AS cu ON ca.customer_id = cu.customer_id
		LEFT JOIN products AS pr ON ca.product_id = pr.product_id
		WHERE %s
		ORDER BY ca.created_at DESC
		%s`, strings.Join(whereQuery, " AND "), limitQuery)

	countQuery := fmt.Sprintf(`SELECT
		COUNT(*)
		FROM cases AS ca
		LEFT JOIN contractors AS co ON ca.contractor_id = co.contractor_id
		LEFT JOIN partners AS pa ON ca.partner_id = pa.partner_id
		LEFT JOIN customers AS cu ON ca.customer_id = cu.customer_id
		LEFT JOIN products AS pr ON ca.product_id = pr.product_id
		WHERE %s`, strings.Join(whereQuery, " AND "))

	var foundCases []CaseFullDTO
	rows, err := r.client.QueryxContext(ctx, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.CaseFull]{}, err
	}
	defer rows.Close()

	caseIDs := make([]string, 0)

	for rows.Next() {
		item := new(CaseFullDTO)
		err := rows.Scan(
			&item.CaseID,
			&item.OwnerID,
			&item.OriginChannel,
			&item.Type,
			&item.Subject,
			&item.Priority,
			&item.Status,
			&item.DueDate,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedBy,
			&item.UpdatedAt,
			&item.ClosedAt,
			&item.TargetDate,
			&item.ExternalReference,
			&item.Region,
			&item.Contractor.ContractorID,
			&item.Contractor.CompanyName,
			&item.Contractor.LegalName,
			&item.Contractor.Document,
			&item.Contractor.DocumentType,
			&item.Contractor.BusinessPhone,
			&item.Contractor.BusinessEmail,
			&item.Contractor.CreatedBy,
			&item.Contractor.CreatedAt,
			&item.Contractor.UpdatedBy,
			&item.Contractor.UpdatedAt,
			&item.Contractor.Active,
			&item.Customer.CustomerID,
			&item.Customer.FirstName,
			&item.Customer.LastName,
			&item.Customer.CompanyName,
			&item.Customer.LegalName,
			&item.Customer.CustomerType,
			&item.Customer.Document,
			&item.Customer.DocumentType,
			&item.Customer.ShippingAddress,
			&item.Customer.ShippingCity,
			&item.Customer.ShippingState,
			&item.Customer.ShippingZipCode,
			&item.Customer.ShippingCountry,
			&item.Customer.BillingAddress,
			&item.Customer.BillingCity,
			&item.Customer.BillingState,
			&item.Customer.BillingZipCode,
			&item.Customer.BillingCountry,
			&item.Customer.PersonalPhone,
			&item.Customer.BusinessPhone,
			&item.Customer.PersonalEmail,
			&item.Customer.BusinessEmail,
			&item.Customer.OwnerID,
			&item.Customer.CreatedBy,
			&item.Customer.CreatedAt,
			&item.Customer.UpdatedBy,
			&item.Customer.UpdatedAt,
			&item.Customer.Active,
			&item.Partner.PartnerID,
			&item.Partner.FirstName,
			&item.Partner.LastName,
			&item.Partner.CompanyName,
			&item.Partner.LegalName,
			&item.Partner.PartnerType,
			&item.Partner.Document,
			&item.Partner.DocumentType,
			&item.Partner.ShippingAddress,
			&item.Partner.ShippingCity,
			&item.Partner.ShippingState,
			&item.Partner.ShippingZipCode,
			&item.Partner.ShippingCountry,
			&item.Partner.BillingAddress,
			&item.Partner.BillingCity,
			&item.Partner.BillingState,
			&item.Partner.BillingZipCode,
			&item.Partner.BillingCountry,
			&item.Partner.PersonalPhone,
			&item.Partner.BusinessPhone,
			&item.Partner.PersonalEmail,
			&item.Partner.BusinessEmail,
			&item.Partner.Region,
			&item.Partner.CreatedBy,
			&item.Partner.CreatedAt,
			&item.Partner.UpdatedBy,
			&item.Partner.UpdatedAt,
			&item.Partner.Active,
			&item.Partner.Description,
			&item.Partner.PaymentKey,
			&item.Partner.PaymentKeyOption,
			&item.Partner.PaymentType,
			&item.Partner.PaymentOwner,
			&item.Partner.PaymentIsSameFromOwner,
			&item.Product.ProductID,
			&item.Product.Name,
			&item.Product.Description,
			&item.Product.Brand,
			&item.Product.Model,
			&item.Product.Value,
			&item.Product.SerialNumber,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
			&item.Product.CreatedBy,
			&item.Product.UpdatedBy,
		)
		if err != nil {
			return domain.PagingResult[domain.CaseFull]{}, err
		}

		caseIDs = append(caseIDs, item.CaseID)

		foundCases = append(foundCases, *item)
	}

	caseTransactions, err := r.getCaseTransactions(ctx, caseIDs)
	if err != nil {
		return domain.PagingResult[domain.CaseFull]{}, err
	}

	for _, t := range caseTransactions {
		caseIdx := slices.IndexFunc(foundCases, func(crmCase CaseFullDTO) bool {
			return crmCase.CaseID == t.CaseID
		})

		if caseIdx > -1 {
			foundCases[caseIdx].Transactions = append(foundCases[caseIdx].Transactions, t)
		}
	}

	var countResult int
	err = r.client.GetContext(ctx, &countResult, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.CaseFull]{}, err
	}

	crmCases := mapCaseFullDTOsToCasesFull(foundCases)

	result := domain.PagingResult[domain.CaseFull]{
		Result: crmCases,
		Paging: domain.Paging{
			Total:  countResult,
			Limit:  filters.Limit,
			Offset: filters.Offset,
		},
	}

	return result, nil
}

func (r *caseRepository) getCaseTransactions(ctx context.Context, caseIDs []string) ([]TransactionDTO, error) {
	if len(caseIDs) == 0 {
		return nil, nil
	}

	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(caseIDs, whereQuery, whereArgs, "case_id")

	query := fmt.Sprintf("SELECT * FROM transactions WHERE %s", strings.Join(whereQuery, " AND "))

	var transactions []TransactionDTO
	err := r.client.SelectContext(ctx, &transactions, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
