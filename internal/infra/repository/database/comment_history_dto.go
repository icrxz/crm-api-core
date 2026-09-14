package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CommentHistoryDTO struct {
	HistoryID string    `db:"history_id"`
	CommentID string    `db:"comment_id"`
	EventName string    `db:"event_name"`
	AuthorID  string    `db:"author_id"`
	OldValues JSONMap   `db:"old_values"`
	NewValues JSONMap   `db:"new_values"`
	CreatedAt time.Time `db:"created_at"`
}

func mapCommentHistoryToCommentHistoryDTO(history domain.CommentHistory) CommentHistoryDTO {
	return CommentHistoryDTO{
		HistoryID: history.HistoryID,
		CommentID: history.CommentID,
		EventName: history.EventName,
		AuthorID:  history.AuthorID,
		OldValues: history.OldValues,
		NewValues: history.NewValues,
		CreatedAt: history.CreatedAt,
	}
}
