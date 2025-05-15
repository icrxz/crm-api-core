package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type attachmentRepository struct {
	db *sqlx.DB
}

func NewAttachmentRepository(db *sqlx.DB) domain.AttachmentRepository {
	return &attachmentRepository{
		db: db,
	}
}

func (r *attachmentRepository) Save(ctx context.Context, attachment domain.Attachment) error {
	return nil
}

func (r *attachmentRepository) SaveBatch(ctx context.Context, attachments []domain.Attachment) error {
	if len(attachments) == 0 {
		return nil
	}

	attachmentsDTO := mapAttachmentsToAttachmentsDTO(attachments)

	_, err := r.db.NamedExecContext(ctx, "INSERT INTO attachments (attachment_id, comment_id, key, file_name, attachment_url, file_extension, size, created_at, created_by) VALUES (:attachment_id, :comment_id, :key, :file_name, :attachment_url, :file_extension, :size, :created_at, :created_by)", attachmentsDTO)
	if err != nil {
		return err
	}

	return nil
}

func (r *attachmentRepository) GetByID(ctx context.Context, attachmentID string) (domain.Attachment, error) {
	if attachmentID == "" {
		return domain.Attachment{}, domain.NewValidationError("attachment_id is required", nil)
	}

	var attachmentDTO AttachmentDTO
	err := r.db.GetContext(ctx, &attachmentDTO, "SELECT * FROM attachments WHERE attachment_id = ?", attachmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Attachment{}, domain.NewNotFoundError("attachment not found", map[string]any{"attachment_id": attachmentID})
		}

		return domain.Attachment{}, err
	}

	return mapAttachmentDTOToAttachment(attachmentDTO), nil
}

func (r *attachmentRepository) GetByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error) {
	if commentID == "" {
		return nil, domain.NewValidationError("comment_id is required", nil)
	}

	var attachmentsDTO []AttachmentDTO
	err := r.db.SelectContext(ctx, &attachmentsDTO, "SELECT * FROM attachments WHERE comment_id = $1", commentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	attachments := mapAttachmentsDTOToAttachments(attachmentsDTO)

	return attachments, nil
}

func (r *attachmentRepository) DeleteManyByComments(ctx context.Context, commentIDs []string) error {
	if len(commentIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In("DELETE FROM attachments WHERE comment_id IN (?)", commentIDs)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}
