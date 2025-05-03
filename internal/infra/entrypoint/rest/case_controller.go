package rest

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type CaseController struct {
	caseService      application.CaseService
	batchCaseService application.BatchCaseService
}

func NewCaseController(
	caseService application.CaseService,
	batchCaseService application.BatchCaseService,
) CaseController {
	return CaseController{
		caseService:      caseService,
		batchCaseService: batchCaseService,
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

func (c *CaseController) SearchCasesFull(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	cases, err := c.caseService.SearchCasesFull(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	searchResult := mapSearchResultToSearchResultDTO(cases, mapCasesFullToCasesFullDTOs)

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

	company := ctx.Request.FormValue("company")
	if company == "" {
		ctx.Error(domain.NewValidationError("company is required", nil))
		return
	}

	fileNameSplit := strings.Split(fileHeader.Filename, ".")
	fileExtension := fileNameSplit[len(fileNameSplit)-1]
	allowExtensions := []string{"csv", "xls", "xlsx"}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.Error(err)
		return
	}
	defer file.Close()

	if !slices.Contains(allowExtensions, fileExtension) {
		ctx.Error(domain.NewValidationError("file must be a csv", nil))
		return
	}

	result, err := c.batchCaseService.CreateBatch(ctx, file, fileHeader.Filename, author, company)
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

	if externalReference := ctx.QueryArray("external_reference"); len(externalReference) > 0 {
		filters.ExternalReference = externalReference
	}

	if state := ctx.QueryArray("state"); len(state) > 0 {
		filters.ShippingState = state
	}

	if startDate := ctx.Query("start_date"); startDate != "" {
		filters.StartDate = &startDate
	}

	if endDate := ctx.Query("end_date"); endDate != "" {
		filters.EndDate = &endDate
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

func (c *CaseController) GetCaseFull(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	caseFull, err := c.caseService.GetCaseFullByID(ctx.Request.Context(), caseID)
	if err != nil {
		ctx.Error(err)
		return
	}

	caseFullDTO := mapCaseFullToCaseFullDTO(*caseFull)

	ctx.JSON(http.StatusOK, caseFullDTO)
}
