package rest

import (
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

	ctx.JSON(201, gin.H{"partner_id": partnerID})
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

	ctx.JSON(200, partnerDTO)
}

func (c *PartnerController) SearchPartners(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	partners, err := c.partnerService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	partnerDTOs := mapPartnersToPartnerDTOs(partners)

	ctx.JSON(200, partnerDTOs)
}

func (c *PartnerController) UpdatePartner(ctx *gin.Context) {
	c.partnerService.Update(ctx.Request.Context(), domain.Partner{})
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

	ctx.JSON(204, nil)
}

func (c *PartnerController) parseQueryToFilters(ctx *gin.Context) domain.PartnerFilters {
	filters := domain.PartnerFilters{}

	if documents := ctx.QueryArray("document"); len(documents) > 0 {
		filters.Document = documents
	}

	if partnerTypes := ctx.QueryArray("partner_type"); len(partnerTypes) > 0 {
		filters.PartnerType = partnerTypes
	}

	if partnerIDs := ctx.QueryArray("partner_id"); len(partnerIDs) > 0 {
		filters.PartnerID = partnerIDs
	}

	if regions := ctx.QueryArray("region"); len(regions) > 0 {
		filters.Region = regions
	}

	return filters
}
