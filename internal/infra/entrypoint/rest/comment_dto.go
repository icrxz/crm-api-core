package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateCommentDTO struct {
	Content     string             `json:"content"`
	CommentType domain.CommentType `json:"comment_type"`
	Attachments []AttachmentDTO    `json:"attachments"`
	CreatedBy   string             `json:"created_by"`
}

type CommentDTO struct {
	CommentID   string    `json:"comment_id"`
	CaseID      string    `json:"case_id"`
	Content     string    `json:"content"`
	CommentType string    `json:"comment_type"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedBy   string    `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func mapCommentDTOToComment(commentDTO CommentDTO) domain.Comment {
	return domain.Comment{
		CommentID:   commentDTO.CommentID,
		CaseID:      commentDTO.CaseID,
		Content:     commentDTO.Content,
		CommentType: domain.CommentType(commentDTO.CommentType),
		CreatedBy:   commentDTO.CreatedBy,
		CreatedAt:   commentDTO.CreatedAt,
		UpdatedBy:   commentDTO.UpdatedBy,
		UpdatedAt:   commentDTO.UpdatedAt,
	}
}

func mapCommentToCommentDTO(comment domain.Comment) CommentDTO {
	return CommentDTO{
		CommentID:   comment.CommentID,
		CaseID:      comment.CaseID,
		Content:     comment.Content,
		CommentType: string(comment.CommentType),
		CreatedBy:   comment.CreatedBy,
		CreatedAt:   comment.CreatedAt,
		UpdatedBy:   comment.UpdatedBy,
		UpdatedAt:   comment.UpdatedAt,
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
	comment, err := domain.NewComment(
		caseID,
		createCommentDTO.Content,
		createCommentDTO.CreatedBy,
		createCommentDTO.CommentType,
	)
	if err != nil {
		return domain.Comment{}, err
	}

	return comment, nil
}
