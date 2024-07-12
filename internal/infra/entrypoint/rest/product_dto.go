package rest

import (
	"github.com/icrxz/crm-api-core/internal/domain"
	"time"
)

type ProductDTO struct {
	ProductID    string    `json:"product_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Value        float64   `json:"value"`
	SerialNumber string    `json:"serial_number"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    string    `json:"created_by"`
	UpdatedBy    string    `json:"updated_by"`
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
