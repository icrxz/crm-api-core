package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"net/http"
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
