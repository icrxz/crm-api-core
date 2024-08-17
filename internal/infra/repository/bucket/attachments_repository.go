package bucket

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type attachmentBucket struct {
	s3Client   *s3.Client
	bucketName string
}

func NewAttachmentBucket(s3Client *s3.Client, bucketName string) domain.AttachmentBucket {
	return &attachmentBucket{
		s3Client:   s3Client,
		bucketName: bucketName,
	}
}

func (b *attachmentBucket) Download(ctx context.Context, fileID string) ([]byte, error) {
	result, err := b.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(b.bucketName),
		Key:    aws.String(fileID),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	file, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, err
	}

	return file, nil
}
