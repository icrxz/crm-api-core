package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"net/http"
)

type CustomerController struct {
	customerService application.CustomerService
}

func NewCustomerController(customerService application.CustomerService) CustomerController {
	return CustomerController{
		customerService: customerService,
	}
}

func (c *CustomerController) CreateCustomer(ctx *gin.Context) {
	var customerDTO *CreateCustomerDTO
	err := ctx.BindJSON(&customerDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	customer, err := mapCreateCustomerDTOToCustomer(*customerDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	customerID, err := c.customerService.Create(ctx.Request.Context(), customer)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"customer_id": customerID})
}

func (c *CustomerController) GetCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.Error(domain.NewValidationError("param customerID cannot be empty", nil))
		return
	}

	customer, err := c.customerService.GetByID(ctx.Request.Context(), customerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	customerDTO := mapCustomerToCustomerDTO(*customer)

	ctx.JSON(200, customerDTO)
}

func (c *CustomerController) SearchCustomers(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	customers, err := c.customerService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	customerDTOs := mapCustomersToCustomerDTOs(customers)

	ctx.JSON(200, customerDTOs)
}

func (c *CustomerController) UpdateCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.Error(domain.NewValidationError("param customerID cannot be empty", nil))
		return
	}

	var updateCustomerDTO *UpdateCustomerDTO
	err := ctx.BindJSON(&updateCustomerDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	updateCustomer := mapUpdateCustomerDTOToUpdateCustomer(*updateCustomerDTO)

	err = c.customerService.Update(ctx.Request.Context(), customerID, updateCustomer)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *CustomerController) DeleteCustomer(ctx *gin.Context) {
	customerID := ctx.Param("customerID")
	if customerID == "" {
		ctx.Error(domain.NewValidationError("param customerID cannot be empty", nil))
		return
	}

	err := c.customerService.Delete(ctx.Request.Context(), customerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(204, nil)
}

func (c *CustomerController) parseQueryToFilters(ctx *gin.Context) domain.CustomerFilters {
	filters := domain.CustomerFilters{}

	if documents := ctx.QueryArray("document"); len(documents) > 0 {
		filters.Document = documents
	}

	if customerIDs := ctx.QueryArray("customer_id"); len(customerIDs) > 0 {
		filters.CustomerID = customerIDs
	}

	if owners := ctx.QueryArray("owner"); len(owners) > 0 {
		filters.OwnerID = owners
	}

	if customerTypes := ctx.QueryArray("type"); len(customerTypes) > 0 {
		filters.CustomerType = customerTypes
	}

	return filters
}
