package application

import (
	"context"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type userService struct {
	userRepository domain.UserRepository
}

//go:generate mockgen -source=user_service.go -destination=mock_application/mock_user_service.go -package=mock_application
type UserService interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	Update(ctx context.Context, userID, author string, user domain.UserUpdate) error
	Delete(ctx context.Context, userID string) error
	Search(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error)
	ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error
}

func NewUserService(userRepository domain.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Create(ctx context.Context, user domain.User) (string, error) {
	return s.userRepository.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, domain.NewValidationError("userID cannot be empty", nil)
	}

	return s.userRepository.GetByID(ctx, userID)
}

func (s *userService) Update(ctx context.Context, userID, author string, userUpdate domain.UserUpdate) error {
	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if userUpdate.Email != nil && *userUpdate.Email != user.Email {
		existingUsers, err := s.userRepository.Search(ctx, domain.UserFilters{
			Email:        []string{*userUpdate.Email},
			PagingFilter: domain.PagingFilter{Limit: 1, Offset: 0},
		})
		if err != nil {
			return err
		}

		for _, existing := range existingUsers.Result {
			if existing.UserID != userID {
				return domain.NewConflictError("email already in use", nil)
			}
		}
	}

	user.MergeUpdate(userUpdate, author)

	return s.userRepository.Update(ctx, *user)
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return domain.NewValidationError("userID cannot be empty", nil)
	}

	return s.userRepository.Delete(ctx, userID)
}

func (s *userService) Search(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error) {
	return s.userRepository.Search(ctx, filters)
}

func (s *userService) ChangePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	if newPassword == "" {
		return domain.NewValidationError("new_password cannot be empty", nil)
	}

	if err := domain.ValidatePasswordComplexity(newPassword); err != nil {
		return err
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if !user.ComparePassword(oldPassword) {
		return domain.NewValidationError("old_password is incorrect", nil)
	}

	if err := user.SetPassword(newPassword); err != nil {
		return err
	}
	user.UpdatedAt = time.Now().UTC()
	user.UpdatedBy = userID

	return s.userRepository.Update(ctx, *user)
}
