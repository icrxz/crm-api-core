package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type AttachmentController struct {
	attachmentService application.AttachmentService
}

func NewAttachmentController(attachmentService application.AttachmentService) AttachmentController {
	return AttachmentController{
		attachmentService: attachmentService,
	}
}

func (c *AttachmentController) Delete(ctx *gin.Context) {
	attachmentID := ctx.Param("attachmentID")
	if attachmentID == "" {
		_ = ctx.Error(domain.NewValidationError("attachmentID is required", nil))
		return
	}

	if err := c.attachmentService.DeleteByID(ctx, attachmentID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
