package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type productRepository struct {
	client *sqlx.DB
}

func NewProductRepository(client *sqlx.DB) domain.ProductRepository {
	return &productRepository{
		client: client,
	}
}

func (r *productRepository) CreateProduct(ctx context.Context, product domain.Product) (string, error) {
	productDTO := mapProductToProductDTO(product)

	_, err := r.client.NamedExecContext(
		ctx,
		"INSERT INTO products "+
			"(product_id, name, description, brand, model, value, serial_number, created_at, updated_at, created_by, updated_by) "+
			"VALUES "+
			"(:product_id, :name, :description, :brand, :model, :value, :serial_number, :created_at, :updated_at, :created_by, :updated_by)",
		productDTO,
	)
	if err != nil {
		return "", err
	}

	return product.ProductID, nil
}

func (r *productRepository) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	if productID == "" {
		return nil, domain.NewValidationError("product_id is required", nil)
	}

	var productDTO ProductDTO
	err := r.client.GetContext(ctx, &productDTO, "SELECT * FROM products WHERE product_id=$1", productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no product found with this id", map[string]any{"product_id": productID})
		}
		return nil, err
	}

	product := mapProductDTOToProduct(productDTO)
	return &product, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, product domain.Product) error {
	productDTO := mapProductToProductDTO(product)

	_, err := r.client.NamedExecContext(
		ctx,
		"UPDATE products "+
			"SET name=:name, description=:description, brand=:brand, model=:model, value=:value, serial_number=:serial_number, updated_at=:updated_at, updated_by=:updated_by "+
			"WHERE product_id=:product_id",
		productDTO,
	)
	if err != nil {
		return err
	}

	return nil
}
