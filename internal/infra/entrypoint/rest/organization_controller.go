package rest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type OrganizationController struct {
	organizationService application.OrganizationService
}

func NewOrganizationController(organizationService application.OrganizationService) OrganizationController {
	return OrganizationController{
		organizationService: organizationService,
	}
}

func (c *OrganizationController) CreateOrganization(ctx *gin.Context) {
	var organizationDTO *CreateOrganizationDTO
	err := ctx.BindJSON(&organizationDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	organization, err := mapCreateOrganizationDTOToOrganization(*organizationDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	organizationID, err := c.organizationService.Create(ctx.Request.Context(), organization)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"organization_id": organizationID})
}

func (c *OrganizationController) UpdateOrganization(ctx *gin.Context) {
	organizationID := ctx.Param("organizationID")
	if organizationID == "" {
		ctx.Error(domain.NewValidationError("param organizationID cannot be empty", nil))
		return
	}

	var updateOrganizationDTO *UpdateOrganizationDTO
	err := ctx.BindJSON(&updateOrganizationDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	updateOrganization := mapUpdateOrganizationDTOToUpdateOrganization(*updateOrganizationDTO)

	if err = c.organizationService.Update(ctx.Request.Context(), organizationID, updateOrganization); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *OrganizationController) GetOrganization(ctx *gin.Context) {
	organizationID := ctx.Param("organizationID")
	if organizationID == "" {
		ctx.Error(domain.NewValidationError("param organizationID cannot be empty", nil))
		return
	}

	organization, err := c.organizationService.GetByID(ctx.Request.Context(), organizationID)
	if err != nil {
		ctx.Error(err)
		return
	}

	organizationDTO := mapOrganizationToOrganizationDTO(*organization)

	ctx.JSON(http.StatusOK, organizationDTO)
}

func (c *OrganizationController) SearchOrganizations(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	organizations, err := c.organizationService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	organizationResponse := mapSearchResultToSearchResultDTO(organizations, mapOrganizationsToOrganizationDTOs)

	ctx.JSON(http.StatusOK, organizationResponse)
}

func (c *OrganizationController) DeleteOrganization(ctx *gin.Context) {
	organizationID := ctx.Param("organizationID")
	if organizationID == "" {
		ctx.Error(domain.NewValidationError("param organizationID cannot be empty", nil))
		return
	}

	err := c.organizationService.Delete(ctx.Request.Context(), organizationID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *OrganizationController) parseQueryToFilters(ctx *gin.Context) domain.OrganizationFilters {
	filters := domain.OrganizationFilters{
		PagingFilter: domain.PagingFilter{
			Limit:  10,
			Offset: 0,
		},
	}

	if organizationIDs := ctx.QueryArray("organization_id"); len(organizationIDs) > 0 {
		filters.OrganizationID = organizationIDs
	}

	if names := ctx.QueryArray("name"); len(names) > 0 {
		filters.Name = names
	}

	if active := ctx.Query("active"); active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err != nil {
			return filters
		}
		filters.Active = &activeBool
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		parsedLimit, err := strconv.Atoi(limitParam)
		if err == nil {
			filters.Limit = parsedLimit
		}
	}

	if offsetParam := ctx.Query("offset"); offsetParam != "" {
		parsedOffset, err := strconv.Atoi(offsetParam)
		if err == nil {
			filters.Offset = parsedOffset
		}
	}

	return filters
}
