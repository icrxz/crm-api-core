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
	GetByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, filters domain.UserFilters) ([]domain.User, error)
}

func NewUserService(userRepository domain.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (us *userService) Create(ctx context.Context, user domain.User) (string, error) {
	return us.userRepository.Create(ctx, user)
}

func (us *userService) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return us.userRepository.GetByID(ctx, id)
}

func (us *userService) Update(ctx context.Context, user domain.User) error {
	return us.userRepository.Update(ctx, user)
}

func (us *userService) Delete(ctx context.Context, id string) error {
	return us.userRepository.Delete(ctx, id)
}

func (us *userService) Search(ctx context.Context, filters domain.UserFilters) ([]domain.User, error) {
	return us.userRepository.Search(ctx, filters)
}
