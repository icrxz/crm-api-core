package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type productService struct {
	productRepository domain.ProductRepository
}

type ProductService interface {
	CreateProduct(ctx context.Context, product domain.Product) (string, error)
	GetProductByID(ctx context.Context, productID string) (*domain.Product, error)
}

func NewProductService(productRepository domain.ProductRepository) ProductService {
	return &productService{
		productRepository: productRepository,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product domain.Product) (string, error) {
	return s.productRepository.CreateProduct(ctx, product)
}

func (s *productService) GetProductByID(ctx context.Context, productID string) (*domain.Product, error) {
	if productID == "" {
		return nil, domain.NewValidationError("productID is required", nil)
	}

	return s.productRepository.GetProductByID(ctx, productID)
}
