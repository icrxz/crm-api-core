package application

import (
	"context"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type authService struct {
	userRepository domain.UserRepository
}

type AuthService interface {
	Login(ctx context.Context) error
	Logout(ctx context.Context) error
}

func NewAuthService(userRepository domain.UserRepository) AuthService {
	return &authService{
		userRepository: userRepository,
	}
}

// Login implements AuthService.
func (a *authService) Login(ctx context.Context) error {
	panic("unimplemented")
}

// Logout implements AuthService.
func (a *authService) Logout(ctx context.Context) error {
	panic("unimplemented")
}
