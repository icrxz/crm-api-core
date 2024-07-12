package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type ProductDTO struct {
	ProductID    string    `db:"product_id"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	Brand        string    `db:"brand"`
	Model        string    `db:"model"`
	Value        float64   `db:"value"`
	SerialNumber string    `db:"serial_number"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
	CreatedBy    string    `db:"created_by"`
	UpdatedBy    string    `db:"updated_by"`
}

func mapProductToProductDTO(product domain.Product) ProductDTO {
	return ProductDTO{
		ProductID:    product.ProductID,
		Name:         product.Name,
		Description:  product.Description,
		Brand:        product.Brand,
		Model:        product.Model,
		Value:        product.Value,
		SerialNumber: product.SerialNumber,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		CreatedBy:    product.CreatedBy,
		UpdatedBy:    product.UpdatedBy,
	}
}

func mapProductDTOToProduct(productDTO ProductDTO) domain.Product {
	return domain.Product{
		ProductID:    productDTO.ProductID,
		Name:         productDTO.Name,
		Description:  productDTO.Description,
		Brand:        productDTO.Brand,
		Model:        productDTO.Model,
		Value:        productDTO.Value,
		SerialNumber: productDTO.SerialNumber,
		CreatedAt:    productDTO.CreatedAt,
		UpdatedAt:    productDTO.UpdatedAt,
		CreatedBy:    productDTO.CreatedBy,
		UpdatedBy:    productDTO.UpdatedBy,
	}
}
