package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type CommentController struct {
	commentService application.CommentService
	userService    application.UserService
}

func NewCommentController(commentService application.CommentService, userService application.UserService) CommentController {
	return CommentController{
		commentService: commentService,
		userService:    userService,
	}
}

func (c *CommentController) canEditComment(ctx *gin.Context, comment *domain.Comment) (string, error) {
	requesterID := ctx.GetString("user_id")

	requester, err := c.userService.GetByID(ctx.Request.Context(), requesterID)
	if err != nil {
		return "", err
	}

	if !requester.Role.IsAdmin() && requesterID != comment.CreatedBy {
		return "", domain.NewUnauthorizedError("user cannot edit another user's comment")
	}

	return requesterID, nil
}

func (c *CommentController) CreateComment(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("caseID is required", nil))
		return
	}

	var commentDTO CreateCommentDTO
	if err := ctx.ShouldBindJSON(&commentDTO); err != nil {
		ctx.Error(domain.NewValidationError("invalid request body", nil))
		return
	}

	comment, err := mapCreateCommentDTOToComment(commentDTO, caseID)
	if err != nil {
		ctx.Error(err)
		return
	}

	commentID, err := c.commentService.Create(ctx, comment)
	if err != nil {
		ctx.Error(err)
		return
	}
	comment.CommentID = commentID

	ctx.JSON(http.StatusCreated, mapCommentToCommentDTO(comment))
}

func (c *CommentController) AddAttachment(ctx *gin.Context) {
	commentID := ctx.Param("commentID")
	if commentID == "" {
		_ = ctx.Error(domain.NewValidationError("commentID is required", nil))
		return
	}

	comment, err := c.commentService.GetByID(ctx, commentID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if _, err := c.canEditComment(ctx, comment); err != nil {
		_ = ctx.Error(err)
		return
	}

	var attachmentDTO CreateAttachmentDTO
	if err := ctx.ShouldBindJSON(&attachmentDTO); err != nil {
		_ = ctx.Error(domain.NewValidationError("invalid request body", nil))
		return
	}

	attachment, err := mapCreateAttachmentDTOToAttachment(attachmentDTO)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	savedAttachment, err := c.commentService.AddAttachment(ctx, commentID, attachment)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, mapAttachmentToAttachmentDTO(savedAttachment))
}

func (c *CommentController) UpdateContent(ctx *gin.Context) {
	commentID := ctx.Param("commentID")
	if commentID == "" {
		_ = ctx.Error(domain.NewValidationError("commentID is required", nil))
		return
	}

	var updateDTO UpdateCommentDTO
	if err := ctx.ShouldBindJSON(&updateDTO); err != nil {
		_ = ctx.Error(domain.NewValidationError("invalid request body", nil))
		return
	}

	comment, err := c.commentService.GetByID(ctx, commentID)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	requesterID, err := c.canEditComment(ctx, comment)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	if err := c.commentService.UpdateContent(ctx, commentID, updateDTO.Content, requesterID); err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *CommentController) GetByID(ctx *gin.Context) {
	commentID := ctx.Param("commentID")
	if commentID == "" {
		ctx.Error(domain.NewValidationError("commentID is required", nil))
		return
	}

	comment, err := c.commentService.GetByID(ctx, commentID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapCommentToCommentDTO(*comment))
}

func (c *CommentController) GetByCaseID(ctx *gin.Context) {
	caseID := ctx.Param("caseID")
	if caseID == "" {
		ctx.Error(domain.NewValidationError("caseID is required", nil))
		return
	}

	comments, err := c.commentService.GetByCaseID(ctx, caseID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mapCommentsToCommentDTOs(comments))
}
