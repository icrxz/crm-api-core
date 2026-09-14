package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=comment_history.go -destination=mock_domain/mock_comment_history_repository.go -package=mock_domain
type CommentHistoryRepository interface {
	Create(ctx context.Context, history CommentHistory) error
}

type CommentHistory struct {
	HistoryID string
	CommentID string
	EventName string
	AuthorID  string
	OldValues map[string]any
	NewValues map[string]any
	CreatedAt time.Time
}

const (
	CommentContentUpdatedEvent    = "comment_content_updated"
	CommentAttachmentAddedEvent   = "comment_attachment_added"
	CommentAttachmentRemovedEvent = "comment_attachment_removed"
)

func NewCommentHistory(
	commentID string,
	eventName string,
	authorID string,
	oldValues map[string]any,
	newValues map[string]any,
) (CommentHistory, error) {
	historyID, err := uuid.NewUUID()
	if err != nil {
		return CommentHistory{}, err
	}

	return CommentHistory{
		HistoryID: historyID.String(),
		CommentID: commentID,
		EventName: eventName,
		AuthorID:  authorID,
		OldValues: oldValues,
		NewValues: newValues,
		CreatedAt: time.Now().UTC(),
	}, nil
}
