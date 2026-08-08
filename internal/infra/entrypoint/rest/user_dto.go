package rest

import (
	"encoding/json"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

// NullableTime distinguishes an absent JSON field (no change) from a JSON
// field explicitly sent as null (clear the value) or a valid time value.
// UnmarshalJSON only runs when the key is present in the payload, so
// Set is true iff the client sent the field at all.
type NullableTime struct {
	Set   bool
	Value *time.Time
}

func (n NullableTime) MarshalJSON() ([]byte, error) {
	if n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}

func (n *NullableTime) UnmarshalJSON(data []byte) error {
	n.Set = true

	if string(data) == "null" {
		n.Value = nil
		return nil
	}

	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	n.Value = &t
	return nil
}

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

type UpdateUserDTO struct {
	FirstName     *string          `json:"first_name"`
	LastName      *string          `json:"last_name"`
	Email         *string          `json:"email"`
	Role          *domain.UserRole `json:"role"`
	Region        *int             `json:"region"`
	Active        *bool            `json:"active"`
	LastAbsenceAt NullableTime     `json:"last_absence_at"`
	UpdatedBy     string           `json:"created_by"`
}

type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type UserDTO struct {
	UserID        string          `json:"user_id"`
	Username      string          `json:"username"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	Email         string          `json:"email"`
	Role          domain.UserRole `json:"role"`
	Region        int             `json:"region"`
	CreatedAt     time.Time       `json:"created_at"`
	CreatedBy     string          `json:"created_by"`
	UpdatedAt     time.Time       `json:"updated_at"`
	UpdatedBy     string          `json:"updated_by"`
	Active        bool            `json:"active"`
	LastAbsenceAt *time.Time      `json:"last_absence_at"`
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
		UserID:        user.UserID,
		Username:      user.Username,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Email:         user.Email,
		Role:          user.Role,
		CreatedAt:     user.CreatedAt,
		CreatedBy:     user.CreatedBy,
		UpdatedAt:     user.UpdatedAt,
		UpdatedBy:     user.UpdatedBy,
		Region:        user.Region,
		Active:        user.Active,
		LastAbsenceAt: user.LastAbsenceAt,
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

func mapUpdateUserDTOToUserUpdate(dto UpdateUserDTO) domain.UserUpdate {
	return domain.UserUpdate{
		FirstName:     dto.FirstName,
		LastName:      dto.LastName,
		Email:         dto.Email,
		Role:          dto.Role,
		Region:        dto.Region,
		Active:        dto.Active,
		LastAbsenceAt: domain.OptionalTime{Present: dto.LastAbsenceAt.Set, Value: dto.LastAbsenceAt.Value},
	}
}
