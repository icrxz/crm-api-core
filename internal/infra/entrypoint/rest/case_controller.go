package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"net/http"
)

type CaseController struct {
	caseService application.CaseService
}

func NewCaseController(
	caseService application.CaseService,
) CaseController {
	return CaseController{
		caseService: caseService,
	}
}

func (c *CaseController) CreateCase(ctx *gin.Context) {
	var createCaseDTO *CreateCaseDTO
	if err := ctx.BindJSON(&createCaseDTO); err != nil {
		ctx.Error(err)
		return
	}

	newCase, err := mapCreateCaseDTOToCreateCase(*createCaseDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	caseID, err := c.caseService.CreateCase(ctx.Request.Context(), newCase)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"case_id": caseID})
}

func (c *CaseController) GetCase(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	crmCase, err := c.caseService.GetCaseByID(ctx.Request.Context(), caseID)
	if err != nil {
		ctx.Error(err)
		return
	}

	caseDTO := mapCaseToCaseDTO(*crmCase)

	ctx.JSON(http.StatusOK, caseDTO)
}

func (c *CaseController) SearchCases(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	cases, err := c.caseService.SearchCases(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	caseDTOs := mapCasesToCaseDTOs(cases)

	ctx.JSON(http.StatusOK, caseDTOs)
}

func (c *CaseController) parseQueryToFilters(ctx *gin.Context) domain.CaseFilters {
	filters := domain.CaseFilters{}

	if ownerIDs := ctx.QueryArray("owner_id"); len(ownerIDs) > 0 {
		filters.OwnerID = ownerIDs
	}

	if contractorIDs := ctx.QueryArray("contractor_id"); len(contractorIDs) > 0 {
		filters.ContractorID = contractorIDs
	}

	if partnerIDs := ctx.QueryArray("partner_id"); len(partnerIDs) > 0 {
		filters.PartnerID = partnerIDs
	}

	if customerIDs := ctx.QueryArray("customer_id"); len(customerIDs) > 0 {
		filters.CustomerID = customerIDs
	}

	if status := ctx.QueryArray("status"); len(status) > 0 {
		filters.Status = status
	}

	if region := ctx.QueryArray("region"); len(region) > 0 {
		filters.Region = region
	}

	return filters
}
