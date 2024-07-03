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

type partnerRepository struct {
	client *sqlx.DB
}

func NewPartnerRepository(client *sqlx.DB) domain.PartnerRepository {
	return &partnerRepository{
		client: client,
	}
}

func (db *partnerRepository) Create(ctx context.Context, partner domain.Partner) (string, error) {
	partnerDTO := mapPartnerToPartnerDTO(partner)

	_, err := db.client.NamedExecContext(
		ctx,
		"INSERT INTO partners "+
			"(partner_id, first_name, last_name, company_name, legal_name, partner_type, document, document_type, shipping_address, shipping_city, shipping_state, shipping_zip_code, shipping_country, billing_address, billing_city, billing_state, billing_zip_code, billing_country, personal_phone, business_phone, personal_email, business_email, created_at, created_by, updated_at, updated_by, active) "+
			"VALUES "+
			"(:partner_id, :first_name, :last_name, :company_name, :legal_name, :partner_type, :document, :document_type, :shipping_address, :shipping_city, :shipping_state, :shipping_zip_code, :shipping_country, :billing_address, :billing_city, :billing_state, :billing_zip_code, :billing_country, :personal_phone, :business_phone, :personal_email, :business_email, :created_at, :created_by, :updated_at, :updated_by, :active)",
		partnerDTO,
	)
	if err != nil {
		return "", err
	}

	return partner.PartnerID, nil
}

func (db *partnerRepository) Delete(ctx context.Context, partnerID string) error {
	if partnerID == "" {
		return domain.NewValidationError("partnerID is required", map[string]any{"partner_id": partnerID})
	}

	_, err := db.client.ExecContext(ctx, "UPDATE partners SET active = false WHERE partner_id = $1", partnerID)
	if err != nil {
		return err
	}

	return nil
}

func (db *partnerRepository) GetByID(ctx context.Context, partnerID string) (*domain.Partner, error) {
	var partnerDTO PartnerDTO
	err := db.client.GetContext(ctx, &partnerDTO, "SELECT * FROM partners WHERE partner_id=$1", partnerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no partner found with this id", map[string]any{"partner_id": partnerID})
		}
		return nil, err
	}

	partner := mapPartnerDTOToPartner(partnerDTO)

	return &partner, nil
}

func (db *partnerRepository) Search(ctx context.Context, filters domain.PartnerFilters) ([]domain.Partner, error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.Document, whereQuery, whereArgs, "document")
	whereQuery, whereArgs = prepareInQuery(filters.PartnerType, whereQuery, whereArgs, "partner_type")
	whereQuery, whereArgs = prepareInQuery(filters.PartnerID, whereQuery, whereArgs, "partner_id")
	whereQuery, whereArgs = prepareInQuery(filters.State, whereQuery, whereArgs, "shipping_state")
	if filters.Active != nil {
		whereQuery = append(whereQuery, fmt.Sprintf("active = $%d", len(whereArgs)+1))
		whereArgs = append(whereArgs, filters.Active)

	}

	query := fmt.Sprintf("SELECT * FROM partners WHERE %s", strings.Join(whereQuery, " AND "))

	var foundPartners []PartnerDTO
	err := db.client.SelectContext(ctx, &foundPartners, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	partners := mapPartnerDTOsToPartners(foundPartners)

	return partners, nil
}

func (db *partnerRepository) Update(ctx context.Context, partner domain.Partner) error {
	partnerDTO := mapPartnerToPartnerDTO(partner)

	_, err := db.client.NamedExecContext(
		ctx,
		"UPDATE partners "+
			"SET first_name = :first_name, "+
			"last_name = :last_name, "+
			"company_name = :company_name, "+
			"legal_name = :legal_name, "+
			"partner_type = :partner_type, "+
			"document = :document, "+
			"document_type = :document_type, "+
			"shipping_address = :shipping_address, "+
			"shipping_city = :shipping_city, "+
			"shipping_state = :shipping_state, "+
			"shipping_zip_code = :shipping_zip_code, "+
			"shipping_country = :shipping_country, "+
			"billing_address = :billing_address, "+
			"billing_city = :billing_city, "+
			"billing_state = :billing_state, "+
			"billing_zip_code = :billing_zip_code, "+
			"billing_country = :billing_country, "+
			"personal_phone = :personal_phone, "+
			"business_phone = :business_phone, "+
			"personal_email = :personal_email, "+
			"business_email = :business_email, "+
			"updated_at = :updated_at, "+
			"updated_by = :updated_by, "+
			"active = :active "+
			"WHERE partner_id = :partner_id",
		partnerDTO,
	)
	if err != nil {
		return err
	}

	return nil
}
