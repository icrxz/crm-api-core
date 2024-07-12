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

type customerRepository struct {
	client *sqlx.DB
}

func NewCustomerRepository(client *sqlx.DB) domain.CustomerRepository {
	return &customerRepository{
		client: client,
	}
}

func (db *customerRepository) Create(ctx context.Context, customer domain.Customer) (string, error) {
	customerDTO := mapCustomerToCustomerDTO(customer)

	_, err := db.client.NamedExecContext(
		ctx,
		"INSERT INTO customers "+
			"(customer_id, first_name, last_name, company_name, legal_name, customer_type, document, document_type, shipping_address, shipping_city, shipping_state, shipping_zip_code, shipping_country, billing_address, billing_city, billing_state, billing_zip_code, billing_country, personal_phone, business_phone, personal_email, business_email, created_at, created_by, updated_at, updated_by, active) "+
			"VALUES "+
			"(:customer_id, :first_name, :last_name, :company_name, :legal_name, :customer_type, :document, :document_type, :shipping_address, :shipping_city, :shipping_state, :shipping_zip_code, :shipping_country, :billing_address, :billing_city, :billing_state, :billing_zip_code, :billing_country, :personal_phone, :business_phone, :personal_email, :business_email, :created_at, :created_by, :updated_at, :updated_by, :active)",
		customerDTO,
	)
	if err != nil {
		return "", err
	}

	return customer.CustomerID, nil
}

func (db *customerRepository) GetByID(ctx context.Context, customerID string) (*domain.Customer, error) {
	var customerDTO CustomerDTO
	err := db.client.GetContext(ctx, &customerDTO, "SELECT * FROM customers WHERE customer_id=$1", customerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no customer found with this id", map[string]any{"customer_id": customerID})
		}
		return nil, err
	}

	customer := mapCustomerDTOToCustomer(customerDTO)

	return &customer, nil
}

func (db *customerRepository) Search(ctx context.Context, filters domain.CustomerFilters) ([]domain.Customer, error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.Document, whereQuery, whereArgs, "document")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerType, whereQuery, whereArgs, "customer_type")
	whereQuery, whereArgs = prepareInQuery(filters.CustomerID, whereQuery, whereArgs, "customer_id")

	query := fmt.Sprintf("SELECT * FROM customers WHERE %s", strings.Join(whereQuery, " AND "))

	var foundCustomers []CustomerDTO
	err := db.client.SelectContext(ctx, &foundCustomers, query, whereArgs...)
	if err != nil {
		return nil, err
	}

	customers := mapCustomerDTOsToCustomers(foundCustomers)

	return customers, nil
}

func (db *customerRepository) Update(ctx context.Context, customer domain.Customer) error {
	customerDTO := mapCustomerToCustomerDTO(customer)

	_, err := db.client.NamedExecContext(
		ctx,
		"UPDATE customers SET "+
			"first_name = :first_name, "+
			"last_name = :last_name, "+
			"company_name = :company_name, "+
			"legal_name = :legal_name, "+
			"customer_type = :customer_type, "+
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
			"WHERE customer_id = :customer_id",
		customerDTO,
	)
	if err != nil {
		return err
	}

	return nil
}

func (db *customerRepository) Delete(ctx context.Context, customerID string) error {
	if customerID == "" {
		return domain.NewValidationError("customer id is required", map[string]any{"customer_id": customerID})
	}

	_, err := db.client.ExecContext(ctx, "UPDATE customers SET active = false WHERE customer_id = $1", customerID)
	if err != nil {
		return err
	}

	return nil
}
