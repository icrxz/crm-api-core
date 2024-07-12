package rest

import "github.com/icrxz/crm-api-core/internal/domain"

type CredentialsDTO struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthUserDTO struct {
	UserID string `json:"user_id" validate:"required"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

type AuthResponseDTO struct {
	Token string      `json:"token" validate:"required"`
	User  AuthUserDTO `json:"user" validate:"required"`
}

func mapUserToAuthResponseDTO(token string, user domain.User) AuthResponseDTO {
	return AuthResponseDTO{
		Token: token,
		User: AuthUserDTO{
			UserID: user.UserID,
			Name:   user.FirstName,
			Email:  user.Email,
			Role:   string(user.Role),
		},
	}
}
