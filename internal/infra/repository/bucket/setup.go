package bucket

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/icrxz/crm-api-core/internal/infra/config"
)

func NewS3Bucket(ctx context.Context, bucketConfig config.Bucket) (*s3.Client, error) {
	sdkConfig, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(bucketConfig.Region))
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(sdkConfig)

	return s3Client, nil
}
