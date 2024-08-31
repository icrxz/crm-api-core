package rest

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
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

	searchResult := mapSearchResultToSearchResultDTO(cases, mapCasesToCaseDTOs)

	ctx.JSON(http.StatusOK, searchResult)
}

func (c *CaseController) CreateBatch(ctx *gin.Context) {
	author := ctx.GetHeader("X-Author")
	if author == "" {
		ctx.Error(domain.NewValidationError("header X-Author cannot be empty", nil))
		return
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.Error(err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.Error(err)
		return
	}
	defer file.Close()

	if !strings.Contains(fileHeader.Filename, ".csv") {
		ctx.Error(domain.NewValidationError("file must be a csv", nil))
		return
	}

	result, err := c.caseService.CreateBatch(ctx, file, author)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"case_ids": result})
}

func (c *CaseController) parseQueryToFilters(ctx *gin.Context) domain.CaseFilters {
	filters := domain.CaseFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

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

	if limit := ctx.Query("limit"); limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err == nil {
			filters.Limit = parsedLimit
		}
	}

	if offset := ctx.Query("offset"); offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err == nil {
			filters.Offset = parsedOffset
		}
	}

	return filters
}

func (c *CaseController) UpdateCase(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	var updateCaseDTO *UpdateCaseDTO
	if err := ctx.BindJSON(&updateCaseDTO); err != nil {
		ctx.Error(err)
		return
	}

	caseUpdate := mapUpdateCaseDTOToUpdateCase(*updateCaseDTO)

	err := c.caseService.UpdateCase(ctx.Request.Context(), caseID, caseUpdate)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
