package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type UserDTO struct {
	UserID    string    `db:"user_id"`
	Username  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Role      string    `db:"role"`
	Region    int       `db:"region"`
	CreatedAt time.Time `db:"created_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy string    `db:"updated_by"`
}

func mapUserToUserDTO(user domain.User) UserDTO {
	return UserDTO{
		UserID:    user.UserID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      string(user.Role),
		Region:    user.Region,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
	}
}

func mapUserDTOToUser(userDTO UserDTO) domain.User {
	return domain.User{
		UserID:    userDTO.UserID,
		Username:  userDTO.Username,
		FirstName: userDTO.FirstName,
		LastName:  userDTO.LastName,
		Email:     userDTO.Email,
		Role:      domain.UserRole(userDTO.Role),
		Region:    userDTO.Region,
		CreatedAt: userDTO.CreatedAt,
		CreatedBy: userDTO.CreatedBy,
		UpdatedAt: userDTO.UpdatedAt,
		UpdatedBy: userDTO.UpdatedBy,
		Password:  userDTO.Password,
	}
}

func mapUserDTOsToUsers(userDTOs []UserDTO) []domain.User {
	users := make([]domain.User, 0, len(userDTOs))
	for _, userDTO := range userDTOs {
		user := mapUserDTOToUser(userDTO)
		users = append(users, user)
	}

	return users
}
