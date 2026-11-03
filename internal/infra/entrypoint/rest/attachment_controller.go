package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type AttachmentController struct {
	attachmentService application.AttachmentService
	userService       application.UserService
}

func NewAttachmentController(attachmentService application.AttachmentService, userService application.UserService) AttachmentController {
	return AttachmentController{
		attachmentService: attachmentService,
		userService:       userService,
	}
}

func (c *AttachmentController) Delete(ctx *gin.Context) {
	attachmentID := ctx.Param("attachmentID")
	if attachmentID == "" {
		_ = ctx.Error(domain.NewValidationError("attachmentID is required", nil))
		return
	}

	attachment, err := c.attachmentService.GetByID(ctx, attachmentID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	requester, err := c.userService.GetByID(ctx.Request.Context(), ctx.GetString("user_id"))
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	// Attachments store the uploader's username as created_by, not the user id.
	if !requester.Role.IsAdmin() && requester.Username != attachment.CreatedBy {
		_ = ctx.Error(domain.NewForbiddenError("user cannot delete another user's attachment"))
		return
	}

	if err := c.attachmentService.DeleteByID(ctx, attachmentID, requester.Username); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
