package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ProductRepository interface {
	CreateProduct(ctx context.Context, product Product) (string, error)
	UpdateProduct(ctx context.Context, product Product) error
	GetProductByID(ctx context.Context, productID string) (*Product, error)
}

type Product struct {
	ProductID    string
	Name         string
	Description  string
	Value        float64
	Brand        string
	Model        string
	SerialNumber string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CreatedBy    string
	UpdatedBy    string
}

func NewProduct(
	name string,
	description string,
	value float64,
	brand string,
	model string,
	serialNumber string,
	createdBy string,
) (Product, error) {
	now := time.Now().UTC()
	productID, err := uuid.NewUUID()
	if err != nil {
		return Product{}, err
	}

	return Product{
		ProductID:    productID.String(),
		Name:         name,
		Description:  description,
		Value:        value,
		Brand:        brand,
		Model:        model,
		SerialNumber: serialNumber,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedBy:    createdBy,
		UpdatedBy:    createdBy,
	}, nil
}

type UpdateProduct struct {
	Name         *string
	Description  *string
	Value        *float64
	Brand        *string
	Model        *string
	SerialNumber *string
	UpdatedBy    string
}

func (p *Product) MergeUpdate(updateProduct UpdateProduct) {
	p.UpdatedBy = updateProduct.UpdatedBy
	p.UpdatedAt = time.Now().UTC()

	if updateProduct.Name != nil {
		p.Name = *updateProduct.Name
	}

	if updateProduct.Description != nil {
		p.Description = *updateProduct.Description
	}

	if updateProduct.Value != nil {
		p.Value = *updateProduct.Value
	}

	if updateProduct.Brand != nil {
		p.Brand = *updateProduct.Brand
	}

	if updateProduct.Model != nil {
		p.Model = *updateProduct.Model
	}

	if updateProduct.SerialNumber != nil {
		p.SerialNumber = *updateProduct.SerialNumber
	}
}
