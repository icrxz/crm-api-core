package rest

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type CreateUserDTO struct {
	Username  string          `json:"username"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	Region    int             `json:"region"`
	Password  string          `json:"password"`
	CreatedBy string          `json:"created_by"`
}

type UserDTO struct {
	UserID    string          `json:"user_id"`
	Username  string          `json:"username"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     string          `json:"email"`
	Role      domain.UserRole `json:"role"`
	Region    int             `json:"region"`
	CreatedAt time.Time       `json:"created_at"`
	CreatedBy string          `json:"created_by"`
	UpdatedAt time.Time       `json:"updated_at"`
	UpdatedBy string          `json:"updated_by"`
	Active    bool            `json:"active"`
}

func mapCreateUserDTOToUser(userDTO CreateUserDTO) (domain.User, error) {
	user, err := domain.NewUser(
		userDTO.FirstName,
		userDTO.LastName,
		userDTO.Email,
		userDTO.Password,
		userDTO.CreatedBy,
		userDTO.Username,
		userDTO.Role,
		userDTO.Region,
	)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func mapUserToUserDTO(user domain.User) UserDTO {
	return UserDTO{
		UserID:    user.UserID,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		CreatedBy: user.CreatedBy,
		UpdatedAt: user.UpdatedAt,
		UpdatedBy: user.UpdatedBy,
		Region:    user.Region,
		Active:    user.Active,
	}
}

func mapUsersToUserDTOs(users []domain.User) []UserDTO {
	userDTOs := make([]UserDTO, 0, len(users))
	for _, user := range users {
		userDTO := mapUserToUserDTO(user)
		userDTOs = append(userDTOs, userDTO)
	}

	return userDTOs
}
