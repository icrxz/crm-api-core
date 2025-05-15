package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type commentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) domain.CommentRepository {
	return &commentRepository{
		db: db,
	}
}

func (r *commentRepository) Create(ctx context.Context, comment domain.Comment) (string, error) {
	commentDTO := mapCommentToCommentDTO(comment)

	_, err := r.db.NamedExecContext(
		ctx, "INSERT INTO comments "+
			"(comment_id, case_id, content, comment_type, created_by, created_at, updated_by, updated_at) "+
			"VALUES "+
			"(:comment_id, :case_id, :content, :comment_type, :created_by, :created_at, :updated_by, :updated_at)",
		commentDTO,
	)
	if err != nil {
		return "", err
	}

	return comment.CommentID, nil
}

func (r *commentRepository) GetByID(ctx context.Context, commentID string) (*domain.Comment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("commentID is required", nil)
	}

	var commentDTO CommentDTO
	err := r.db.GetContext(ctx, &commentDTO, "SELECT * FROM comments WHERE comment_id = $1", commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no comment found with this id", map[string]any{"comment_id": commentID})
		}
		return nil, err
	}

	comment := mapCommentDTOToComment(commentDTO)

	return &comment, nil
}

func (r *commentRepository) GetByCaseID(ctx context.Context, caseID string) ([]domain.Comment, error) {
	if caseID == "" {
		return nil, domain.NewValidationError("caseID is required", nil)
	}

	var commentDTOs []CommentDTO
	err := r.db.SelectContext(ctx, &commentDTOs, "SELECT * FROM comments WHERE case_id = $1", caseID)
	if err != nil {
		return nil, err
	}

	comments := mapCommentDTOsToComments(commentDTOs)

	return comments, nil
}

func (r *commentRepository) DeleteManyByCaseID(ctx context.Context, caseID string) error {
	if caseID == "" {
		return domain.NewValidationError("caseID is required", nil)
	}

	_, err := r.db.ExecContext(ctx, "DELETE FROM comments WHERE case_id = $1", caseID)
	if err != nil {
		return err
	}

	return nil
}
