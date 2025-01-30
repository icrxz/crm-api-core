package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateProductDTO struct {
	Name         string  `json:"name" binding:"required"`
	Description  string  `json:"description" binding:"required"`
	Value        float64 `json:"value" binding:"required"`
	Brand        string  `json:"brand" binding:"required"`
	Model        string  `json:"model" binding:"required"`
	SerialNumber string  `json:"serial_number" binding:"required"`
	CreatedBy    string  `json:"created_by" binding:"required"`
}

type UpdateProductDTO struct {
	Name         *string  `json:"name"`
	Description  *string  `json:"description"`
	Value        *float64 `json:"value"`
	Brand        *string  `json:"brand"`
	Model        *string  `json:"model"`
	SerialNumber *string  `json:"serial_number"`
	UpdatedBy    string   `json:"updated_by"`
}

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

func mapCreateProductDTOToProduct(productDTO CreateProductDTO) (domain.Product, error) {
	return domain.NewProduct(
		productDTO.Name,
		productDTO.Description,
		productDTO.Value,
		productDTO.Brand,
		productDTO.Model,
		productDTO.SerialNumber,
		productDTO.CreatedBy,
	)
}

func mapUpdateProductDTOToUpdateProduct(productDTO UpdateProductDTO) domain.UpdateProduct {
	return domain.UpdateProduct{
		Name:         productDTO.Name,
		Description:  productDTO.Description,
		Value:        productDTO.Value,
		Brand:        productDTO.Brand,
		Model:        productDTO.Model,
		SerialNumber: productDTO.SerialNumber,
		UpdatedBy:    productDTO.UpdatedBy,
	}
}
