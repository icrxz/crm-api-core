package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type ProductController struct {
	productService application.ProductService
}

func NewProductController(productService application.ProductService) ProductController {
	return ProductController{
		productService: productService,
	}
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
	productID := ctx.Param("productID")
	if productID == "" {
		ctx.Error(domain.NewValidationError("param productID cannot be empty", nil))
		return
	}

	product, err := c.productService.GetProductByID(ctx.Request.Context(), productID)
	if err != nil {
		ctx.Error(err)
		return
	}

	productDTO := mapProductToProductDTO(*product)

	ctx.JSON(http.StatusOK, productDTO)
}

func (c *ProductController) CreateProduct(ctx *gin.Context) {
	var productDTO CreateProductDTO
	if err := ctx.ShouldBindJSON(&productDTO); err != nil {
		ctx.Error(domain.NewValidationError("invalid request body", nil))
		return
	}

	product, err := mapCreateProductDTOToProduct(productDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	productID, err := c.productService.CreateProduct(ctx.Request.Context(), product)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, map[string]string{"productID": productID})
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("productID")
	if productID == "" {
		ctx.Error(domain.NewValidationError("param productID cannot be empty", nil))
		return
	}

	var updateProductDTO *UpdateProductDTO
	if err := ctx.ShouldBindJSON(&updateProductDTO); err != nil {
		ctx.Error(domain.NewValidationError("invalid request body", nil))
		return
	}

	updateProduct := mapUpdateProductDTOToUpdateProduct(*updateProductDTO)

	err := c.productService.UpdateProduct(ctx.Request.Context(), productID, updateProduct)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
