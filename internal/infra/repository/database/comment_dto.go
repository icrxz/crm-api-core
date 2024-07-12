package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CommentDTO struct {
	CommentID   string    `db:"comment_id"`
	CaseID      string    `db:"case_id"`
	Content     string    `db:"content"`
	CommentType string    `db:"comment_type"`
	CreatedBy   string    `db:"created_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedBy   string    `db:"updated_by"`
	UpdatedAt   time.Time `db:"updated_at"`
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

func mapCommentDTOsToComments(commentDTOs []CommentDTO) []domain.Comment {
	var comments []domain.Comment
	for _, commentDTO := range commentDTOs {
		comments = append(comments, mapCommentDTOToComment(commentDTO))
	}
	return comments
}
