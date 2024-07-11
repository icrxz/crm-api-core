package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateCommentDTO struct {
	Content     string                `json:"content"`
	CommentType domain.CommentType    `json:"comment_type"`
	Attachments []CreateAttachmentDTO `json:"attachments"`
	CreatedBy   string                `json:"created_by"`
}

type CommentDTO struct {
	CommentID   string          `json:"comment_id"`
	CaseID      string          `json:"case_id"`
	Content     string          `json:"content"`
	CommentType string          `json:"comment_type"`
	CreatedBy   string          `json:"created_by"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedBy   string          `json:"updated_by"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Attachments []AttachmentDTO `json:"attachments"`
}

func mapCommentToCommentDTO(comment domain.Comment) CommentDTO {
	attachmentsDTOs := mapAttachmentsToAttachmentDTOs(comment.Attachments)

	return CommentDTO{
		CommentID:   comment.CommentID,
		CaseID:      comment.CaseID,
		Content:     comment.Content,
		CommentType: string(comment.CommentType),
		CreatedBy:   comment.CreatedBy,
		CreatedAt:   comment.CreatedAt,
		UpdatedBy:   comment.UpdatedBy,
		UpdatedAt:   comment.UpdatedAt,
		Attachments: attachmentsDTOs,
	}
}

func mapCommentsToCommentDTOs(comments []domain.Comment) []CommentDTO {
	var commentDTOs []CommentDTO
	for _, comment := range comments {
		commentDTOs = append(commentDTOs, mapCommentToCommentDTO(comment))
	}
	return commentDTOs
}

func mapCreateCommentDTOToComment(createCommentDTO CreateCommentDTO, caseID string) (domain.Comment, error) {
	attachmentsDTO, err := mapCreateAttachmentDTOsToAttachments(createCommentDTO.Attachments)
	if err != nil {
		return domain.Comment{}, err
	}

	comment, err := domain.NewComment(
		caseID,
		createCommentDTO.Content,
		createCommentDTO.CreatedBy,
		createCommentDTO.CommentType,
		attachmentsDTO,
	)
	if err != nil {
		return domain.Comment{}, err
	}

	return comment, nil
}
