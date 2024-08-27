package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type userService struct {
	userRepository domain.UserRepository
}

type UserService interface {
	Create(ctx context.Context, user domain.User) (string, error)
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	Update(ctx context.Context, userID, author string, user domain.UserUpdate) error
	Delete(ctx context.Context, userID string) error
	Search(ctx context.Context, filters domain.UserFilters) (domain.PagingResult[domain.User], error)
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
