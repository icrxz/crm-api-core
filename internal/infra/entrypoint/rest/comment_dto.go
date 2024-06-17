package rest

import (
	"github.com/icrxz/crm-api-core/internal/domain"
	"time"
)

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
