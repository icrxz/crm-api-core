package database

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type commentHistoryRepository struct {
	client *sqlx.DB
}

func NewCommentHistoryRepository(client *sqlx.DB) domain.CommentHistoryRepository {
	return &commentHistoryRepository{
		client: client,
	}
}

func (r *commentHistoryRepository) Create(ctx context.Context, history domain.CommentHistory) error {
	historyDTO := mapCommentHistoryToCommentHistoryDTO(history)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		"INSERT INTO comment_history "+
			"(history_id, comment_id, event_name, author_id, old_values, new_values, created_at) "+
			"VALUES "+
			"(:history_id, :comment_id, :event_name, :author_id, :old_values, :new_values, :created_at)",
		historyDTO,
	)

	return err
}
