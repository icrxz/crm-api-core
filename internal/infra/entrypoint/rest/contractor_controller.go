package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
	"net/http"
	"strconv"
)

type ContractorController struct {
	contractorService application.ContractorService
}

func NewContractorController(contractorService application.ContractorService) ContractorController {
	return ContractorController{
		contractorService: contractorService,
	}
}

func (c *ContractorController) CreateContractor(ctx *gin.Context) {
	var contractorDTO *CreateContractorDTO
	err := ctx.BindJSON(&contractorDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	contractor, err := mapCreateContractorDTOToContractor(*contractorDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	contractorID, err := c.contractorService.Create(ctx.Request.Context(), contractor)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(201, gin.H{"contractor_id": contractorID})
}

func (c *ContractorController) UpdateContractor(ctx *gin.Context) {
	contractorID := ctx.Param("contractorID")
	if contractorID == "" {
		ctx.Error(domain.NewValidationError("param contractorID cannot be empty", nil))
		return
	}

	var updateContractorDTO *UpdateContractorDTO
	err := ctx.BindJSON(&updateContractorDTO)
	if err != nil {
		ctx.Error(err)
		return
	}

	updateContractor := mapUpdateContractorDTOToUpdateContractor(*updateContractorDTO)

	if err = c.contractorService.Update(ctx.Request.Context(), contractorID, updateContractor); err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *ContractorController) GetContractor(ctx *gin.Context) {
	contractorID := ctx.Param("contractorID")
	if contractorID == "" {
		ctx.Error(domain.NewValidationError("param contractorID cannot be empty", nil))
		return
	}

	contractor, err := c.contractorService.GetByID(ctx.Request.Context(), contractorID)
	if err != nil {
		ctx.Error(err)
		return
	}

	contractorDTO := mapContractorToContractorDTO(*contractor)

	ctx.JSON(http.StatusOK, contractorDTO)
}

func (c *ContractorController) SearchContractors(ctx *gin.Context) {
	filters := c.parseQueryToFilters(ctx)

	contractors, err := c.contractorService.Search(ctx.Request.Context(), filters)
	if err != nil {
		ctx.Error(err)
		return
	}

	contractorDTOs := mapContractorsToContractorDTOs(contractors)

	ctx.JSON(http.StatusOK, contractorDTOs)
}

func (c *ContractorController) DeleteContractor(ctx *gin.Context) {
	contractorID := ctx.Param("contractorID")
	if contractorID == "" {
		ctx.Error(domain.NewValidationError("param contractorID cannot be empty", nil))
		return
	}

	err := c.contractorService.Delete(ctx.Request.Context(), contractorID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *ContractorController) parseQueryToFilters(ctx *gin.Context) domain.ContractorFilters {
	filters := domain.ContractorFilters{}

	if documents := ctx.QueryArray("document"); len(documents) > 0 {
		filters.Document = documents
	}

	if contractorIDs := ctx.QueryArray("contractor_id"); len(contractorIDs) > 0 {
		filters.ContractorID = contractorIDs
	}

	if companyNames := ctx.QueryArray("company_name"); len(companyNames) > 0 {
		filters.CompanyName = companyNames
	}

	if active := ctx.Query("active"); active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err != nil {
			return filters
		}
		filters.Active = &activeBool
	}

	return filters
}
