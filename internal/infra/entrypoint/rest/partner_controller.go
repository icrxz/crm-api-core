package rest

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type PartnerController struct {
	partnerService application.PartnerService
}

func NewPartnerController(partnerService application.PartnerService) PartnerController {
	return PartnerController{
		partnerService: partnerService,
	}
}

func (c *PartnerController) CreatePartner(ctx *gin.Context) {
	var partnerDTO *CreatePartnerDTO
	err := ctx.BindJSON(&partnerDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	partner, err := mapCreatePartnerDTOToPartner(*partnerDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	partnerID, err := c.partnerService.Create(ctx.Request.Context(), partner)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"partner_id": partnerID})
}

func (c *PartnerController) GetPartner(ctx *gin.Context) {
	partnerID := ctx.Param("partnerID")
	if partnerID == "" {
		ctx.Error(domain.NewValidationError("param partnerID cannot be empty", nil))
		return
	}

	partner, err := c.partnerService.GetByID(ctx.Request.Context(), partnerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	partnerDTO := mapPartnerToPartnerDTO(*partner)

	ctx.JSON(http.StatusOK, partnerDTO)
}

func (c *PartnerController) SearchPartners(ctx *gin.Context) {
	filters, err := c.parseQueryToFilters(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	partners, err := c.partnerService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	searchResult := mapSearchResultToSearchResultDTO(partners, mapPartnersToPartnerDTOs)

	ctx.JSON(http.StatusOK, searchResult)
}

func (c *PartnerController) UpdatePartner(ctx *gin.Context) {
	partnerID := ctx.Param("partnerID")
	if partnerID == "" {
		ctx.Error(domain.NewValidationError("param partnerID cannot be empty", nil))
		return
	}

	var editPartnerDTO *EditPartnerDTO
	err := ctx.BindJSON(&editPartnerDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	editPartner := mapEditPartnerDTOToEditPartner(*editPartnerDTO)

	err = c.partnerService.Update(ctx.Request.Context(), partnerID, editPartner)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *PartnerController) DeletePartner(ctx *gin.Context) {
	partnerID := ctx.Param("partnerID")
	if partnerID == "" {
		ctx.Error(domain.NewValidationError("param partnerID cannot be empty", nil))
		return
	}

	err := c.partnerService.Delete(ctx.Request.Context(), partnerID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *PartnerController) CreateBatch(ctx *gin.Context) {
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

	result, err := c.partnerService.CreateBatch(ctx, file, author)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"partner_ids": result})
}

func (c *PartnerController) parseQueryToFilters(ctx *gin.Context) (domain.PartnerFilters, error) {
	filters := domain.PartnerFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if documents := ctx.QueryArray("document"); len(documents) > 0 {
		filters.Document = documents
	}

	if partnerTypes := ctx.QueryArray("partner_type"); len(partnerTypes) > 0 {
		filters.PartnerType = partnerTypes
	}

	if partnerIDs := ctx.QueryArray("partner_id"); len(partnerIDs) > 0 {
		filters.PartnerID = partnerIDs
	}

	if states := ctx.QueryArray("state"); len(states) > 0 {
		filters.State = states
	}

	if cities := ctx.QueryArray("city"); len(cities) > 0 {
		filters.City = cities
	}

	if firstName := ctx.QueryArray("first_name"); len(firstName) > 0 {
		filters.FirstName = firstName
	}

	if lastName := ctx.QueryArray("last_name"); len(lastName) > 0 {
		filters.LastName = lastName
	}

	if active := ctx.Query("active"); active != "" {
		isActive := active == "true"
		filters.Active = &isActive
	}

	validationErr := make([]error, 0)
	if limitParam := ctx.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err != nil {
			validationErr = append(validationErr, domain.NewValidationError("limit must be a number", nil))
		} else {
			filters.PagingFilter.Limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err != nil {
			validationErr = append(validationErr, domain.NewValidationError("offset must be a number", nil))
		} else {
			filters.PagingFilter.Offset = parsedOffset
		}
	}

	if len(validationErr) > 0 {
		return domain.PartnerFilters{}, errors.Join(validationErr...)
	}

	return filters, nil
}
