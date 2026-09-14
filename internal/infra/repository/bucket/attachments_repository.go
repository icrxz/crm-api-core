package bucket

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type attachmentBucket struct {
	s3Client       *s3.Client
	bucketName     string
	deleteDisabled bool
}

func NewAttachmentBucket(s3Client *s3.Client, bucketName string) domain.AttachmentBucket {
	return &attachmentBucket{
		s3Client:       s3Client,
		bucketName:     bucketName,
		deleteDisabled: os.Getenv("DISABLE_BUCKET_DELETE") == "true",
	}
}

// Delete removes the object from the bucket. Every environment shares the
// same S3 bucket (there is no separate local/dev bucket), so a local
// developer accidentally deleting an attachment would permanently delete the
// production file too. Setting DISABLE_BUCKET_DELETE=true (only meant for
// .env.local, which is gitignored and never present in production) turns
// this into a no-op instead.
func (b *attachmentBucket) Delete(ctx context.Context, key string) error {
	if b.deleteDisabled {
		slog.WarnContext(ctx, "bucket delete skipped because DISABLE_BUCKET_DELETE is set",
			"bucket", b.bucketName, "key", key)
		return nil
	}

	_, err := b.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(key),
	})
	return err
}

func (b *attachmentBucket) Download(ctx context.Context, fileID string) ([]byte, error) {
	result, err := b.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, err
	}
	defer func() { _ = result.Body.Close() }()

	file, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}
