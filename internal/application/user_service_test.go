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

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	user := domain.User{UserID: "user-1", Username: "johndoe"}

	mockRepo.EXPECT().Create(gomock.Any(), user).Return("user-1", nil)

	userID, err := service.Create(context.Background(), user)

	require.NoError(t, err)
	assert.Equal(t, "user-1", userID)
}

func TestUserService_GetByID(t *testing.T) {
	t.Run("returns validation error when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewUserService(mock_domain.NewMockUserRepository(ctrl))

		_, err := service.GetByID(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("returns the user from repository", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		expected := &domain.User{UserID: "user-1", Username: "johndoe"}
		mockRepo.EXPECT().GetByID(gomock.Any(), "user-1").Return(expected, nil)

		user, err := service.GetByID(context.Background(), "user-1")

		require.NoError(t, err)
		assert.Equal(t, expected, user)
	})
}

func TestUserService_Update(t *testing.T) {
	t.Run("merges update and persists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		existing := &domain.User{UserID: "user-1", FirstName: "John"}
		newFirstName := "Jane"

		mockRepo.EXPECT().GetByID(gomock.Any(), "user-1").Return(existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, user domain.User) error {
				assert.Equal(t, newFirstName, user.FirstName)
				return nil
			},
		)

		err := service.Update(context.Background(), "user-1", "author-2", domain.UserUpdate{
			FirstName: &newFirstName,
		})

		require.NoError(t, err)
	})

	t.Run("returns error when user is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		mockRepo.EXPECT().GetByID(gomock.Any(), "user-1").Return(nil, domain.NewNotFoundError("user not found", nil))

		err := service.Update(context.Background(), "user-1", "author-2", domain.UserUpdate{})

		require.Error(t, err)
	})
}

func TestUserService_Delete(t *testing.T) {
	t.Run("returns validation error when userID is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewUserService(mock_domain.NewMockUserRepository(ctrl))

		err := service.Delete(context.Background(), "")

		require.Error(t, err)
	})

	t.Run("deletes the user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		mockRepo.EXPECT().Delete(gomock.Any(), "user-1").Return(nil)

		err := service.Delete(context.Background(), "user-1")

		require.NoError(t, err)
	})
}

func TestUserService_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_domain.NewMockUserRepository(ctrl)
	service := NewUserService(mockRepo)

	filters := domain.UserFilters{Email: []string{"john@doe.com"}}
	expected := domain.PagingResult[domain.User]{
		Result: []domain.User{{UserID: "user-1"}},
	}

	mockRepo.EXPECT().Search(gomock.Any(), filters).Return(expected, nil)

	result, err := service.Search(context.Background(), filters)

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUserService_ChangePassword(t *testing.T) {
	t.Run("returns validation error when new password is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		service := NewUserService(mock_domain.NewMockUserRepository(ctrl))

		err := service.ChangePassword(context.Background(), "user-1", "old-password", "")

		require.Error(t, err)
	})

	t.Run("returns error when user is not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		mockRepo.EXPECT().GetByID(gomock.Any(), "user-1").Return(nil, domain.NewNotFoundError("user not found", nil))

		err := service.ChangePassword(context.Background(), "user-1", "old-password", "new-password")

		require.Error(t, err)
	})

	t.Run("returns validation error when old password does not match", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		existing, err := domain.NewUser("John", "Doe", "john@doe.com", "current-password", "author-1", "johndoe", domain.OPERATOR, 1)
		require.NoError(t, err)

		mockRepo.EXPECT().GetByID(gomock.Any(), "user-1").Return(&existing, nil)

		err = service.ChangePassword(context.Background(), "user-1", "wrong-password", "new-password")

		require.Error(t, err)
	})

	t.Run("changes the password when old password matches", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mock_domain.NewMockUserRepository(ctrl)
		service := NewUserService(mockRepo)

		existing, err := domain.NewUser("John", "Doe", "john@doe.com", "current-password", "author-1", "johndoe", domain.OPERATOR, 1)
		require.NoError(t, err)

		mockRepo.EXPECT().GetByID(gomock.Any(), "user-1").Return(&existing, nil)
		mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ context.Context, user domain.User) error {
				assert.True(t, user.ComparePassword("new-password"))
				return nil
			},
		)

		err = service.ChangePassword(context.Background(), "user-1", "current-password", "new-password")

		require.NoError(t, err)
	})
}
