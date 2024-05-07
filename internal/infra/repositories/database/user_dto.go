package database

import "github.com/icrxz/crm-api-core/internal/domain"

type UserDTO struct {
	UserID   string          `db:"user_id"`
	Name     string          `db:"name"`
	Email    string          `db:"email"`
	Role     domain.UserRole `db:"role"`
	Password string          `db:"password"`
}

func mapUserToUserDTO(user domain.User) UserDTO {
	return UserDTO{
		UserID:   user.UserID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
		Password: user.Password,
	}
}
