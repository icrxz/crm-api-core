package bucket

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAttachmentBucket_Delete_SkipsWhenDisabled(t *testing.T) {
	t.Run("does not touch S3 when deleteDisabled is true", func(t *testing.T) {
		b := &attachmentBucket{deleteDisabled: true}

		err := b.Delete(context.Background(), "some-key.png")

		require.NoError(t, err)
	})

	t.Run("NewAttachmentBucket reads DISABLE_BUCKET_DELETE from the environment", func(t *testing.T) {
		t.Setenv("DISABLE_BUCKET_DELETE", "true")

		b := NewAttachmentBucket(nil, "some-bucket")

		require.NoError(t, b.Delete(context.Background(), "some-key.png"))
	})

	t.Run("NewAttachmentBucket enables real deletes when the env var is unset", func(t *testing.T) {
		bucket := NewAttachmentBucket(nil, "some-bucket").(*attachmentBucket)

		require.False(t, bucket.deleteDisabled)
	})
}
