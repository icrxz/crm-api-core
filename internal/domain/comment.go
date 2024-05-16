package domain

import "time"

type Attachment struct {
	AttachmentID  string
	FileName      string
	AttachmentURL string
	FileExtension string
	CreatedAt     time.Time
}

type Comment struct {
	CommentID   string
	CaseID      string
	Content     string
	CommentType CommentType
	Attachment  []Attachment
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

type CommentType string

const (
	CONTENT    CommentType = "Content"
	COMMENT    CommentType = "Comment"
	RESOLUTION CommentType = "Resolution"
	REJECTION  CommentType = "Rejection"
)
