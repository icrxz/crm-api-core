package application

import (
	"context"
	"testing"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/icrxz/crm-api-core/internal/domain/mock_domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestQueueService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockQueueRepository(ctrl)
	service := NewQueueService(mockRepo)

	queue := domain.Queue{QueueID: "queue-1", Name: "SP Mobile"}

	mockRepo.EXPECT().Create(gomock.Any(), queue).Return("queue-1", nil)

	queueID, err := service.Create(context.Background(), queue)

	require.NoError(t, err)
	assert.Equal(t, "queue-1", queueID)
}

func TestQueueService_GetByID(t *testing.T) {
	t.Run("returns validation error when queueID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		_, err := service.GetByID(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("returns the queue from repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		expected := &domain.Queue{QueueID: "queue-1", Name: "SP Mobile"}
		mockRepo.EXPECT().GetByID(gomock.Any(), "queue-1").Return(expected, nil)

		queue, err := service.GetByID(context.Background(), "queue-1")

		require.NoError(t, err)
		assert.Equal(t, expected, queue)
	})
}

func TestQueueService_Update(t *testing.T) {
	t.Run("returns validation error when queueID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		err := service.Update(context.Background(), "", domain.UpdateQueue{})

		require.Error(t, err)
	})

	t.Run("merges update and persists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		existing := &domain.Queue{QueueID: "queue-1", Name: "SP Mobile", Active: true}
		newName := "SP Mobile Updated"

		mockRepo.EXPECT().GetByID(gomock.Any(), "queue-1").Return(existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, queue domain.Queue) error {
				assert.Equal(t, newName, queue.Name)
				return nil
			},
		)

		err := service.Update(context.Background(), "queue-1", domain.UpdateQueue{
			Name:      &newName,
			UpdatedBy: "author-2",
		})

		require.NoError(t, err)
	})
}

func TestQueueService_Delete(t *testing.T) {
	t.Run("returns validation error when queueID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		err := service.Delete(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("deletes the queue", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		mockRepo.EXPECT().Delete(gomock.Any(), "queue-1").Return(nil)

		err := service.Delete(context.Background(), "queue-1")

		require.NoError(t, err)
	})
}

func TestQueueService_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockQueueRepository(ctrl)
	service := NewQueueService(mockRepo)

	filters := domain.QueueFilters{Category: []string{"mobile"}}
	expected := domain.PagingResult[domain.Queue]{
		Result: []domain.Queue{{QueueID: "queue-1"}},
	}

	mockRepo.EXPECT().Search(gomock.Any(), filters).Return(expected, nil)

	result, err := service.Search(context.Background(), filters)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestQueueService_AddMember(t *testing.T) {
	t.Run("returns validation error when queueID or userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		err := service.AddMember(context.Background(), "", "user-1")
		require.Error(t, err)

		err = service.AddMember(context.Background(), "queue-1", "")
		require.Error(t, err)
	})

	t.Run("adds the member", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		mockRepo.EXPECT().AddMember(gomock.Any(), "queue-1", "user-1").Return(nil)

		err := service.AddMember(context.Background(), "queue-1", "user-1")

		require.NoError(t, err)
	})
}

func TestQueueService_RemoveMember(t *testing.T) {
	t.Run("returns validation error when queueID or userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		err := service.RemoveMember(context.Background(), "", "user-1")
		require.Error(t, err)

		err = service.RemoveMember(context.Background(), "queue-1", "")
		require.Error(t, err)
	})

	t.Run("removes the member", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		mockRepo.EXPECT().RemoveMember(gomock.Any(), "queue-1", "user-1").Return(nil)

		err := service.RemoveMember(context.Background(), "queue-1", "user-1")

		require.NoError(t, err)
	})
}

func TestQueueService_GetMembers(t *testing.T) {
	t.Run("returns validation error when queueID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		_, err := service.GetMembers(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("returns the members", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		expected := []domain.User{{UserID: "user-1"}}
		mockRepo.EXPECT().GetMembers(gomock.Any(), "queue-1").Return(expected, nil)

		members, err := service.GetMembers(context.Background(), "queue-1")

		require.NoError(t, err)
		assert.Equal(t, expected, members)
	})
}

func TestQueueService_GetQueuesByUser(t *testing.T) {
	t.Run("returns validation error when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewQueueService(mock_domain.NewMockQueueRepository(ctrl))

		_, err := service.GetQueuesByUser(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("returns the queues for the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockQueueRepository(ctrl)
		service := NewQueueService(mockRepo)

		expected := []domain.Queue{{QueueID: "queue-1"}}
		mockRepo.EXPECT().GetQueuesByUser(gomock.Any(), "user-1").Return(expected, nil)

		queues, err := service.GetQueuesByUser(context.Background(), "user-1")

		require.NoError(t, err)
		assert.Equal(t, expected, queues)
	})
}
