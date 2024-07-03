package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type CaseActionController struct {
	caseActionService application.CaseActionService
}

func NewCaseActionController(
	caseActionService application.CaseActionService,
) CaseActionController {
	return CaseActionController{
		caseActionService: caseActionService,
	}
}

func (c *CaseActionController) ChangeOwner(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	var changeOwnerDTO ChangeOwnerDTO
	if err := ctx.BindJSON(&changeOwnerDTO); err != nil {
		ctx.Error(err)
		return
	}

	changeOwner := mapChangeOwnerDTOToChangeOwner(changeOwnerDTO)

	err := c.caseActionService.ChangeOwner(ctx.Request.Context(), caseID, changeOwner)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *CaseActionController) ChangeStatus(ctx *gin.Context) {
	fmt.Println("entrou na rota")
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	var changeStatusDTO ChangeStatusDTO
	if err := ctx.BindJSON(&changeStatusDTO); err != nil {
		ctx.Error(err)
		return
	}

	changeStatus := mapChangeStatusDTOToChangeStatus(changeStatusDTO)

	err := c.caseActionService.ChangeStatus(ctx.Request.Context(), caseID, changeStatus)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *CaseActionController) ChangePartner(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	var changePartnerDTO ChangePartnerDTO
	if err := ctx.BindJSON(&changePartnerDTO); err != nil {
		ctx.Error(err)
		return
	}

	changePartner := mapChangePartnerDTOToChangePartner(changePartnerDTO)

	err := c.caseActionService.ChangePartner(ctx.Request.Context(), caseID, changePartner)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (c *CaseActionController) DownloadReport(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("case_id is required", nil))
		return
	}

	report, filename, err := c.caseActionService.GenerateReport(ctx.Request.Context(), caseID)
	if err != nil {
		ctx.Error(err)
		return
	}

	contentType := fmt.Sprintf("application/vnd.openxmlformats-officedocument.wordprocessingml.document;%s.docx", filename)
	ctx.Data(http.StatusOK, contentType, report)
}
