package database

import (
	"testing"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapUserToUserDTO(t *testing.T) {
	t.Run("maps a user with last absence at set", func(t *testing.T) {
		now := time.Now().UTC()
		absenceAt := now.Add(-time.Hour)
		user := domain.User{
			UserID:        "user-1",
			Username:      "johndoe",
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "john@doe.com",
			Role:          domain.OPERATOR,
			Region:        1,
			Password:      "hash",
			LastLoggedIP:  "127.0.0.1",
			SessionToken:  "token",
			Active:        true,
			LastAbsenceAt: &absenceAt,
			CreatedBy:     "author-1",
			CreatedAt:     now,
			UpdatedBy:     "author-1",
			UpdatedAt:     now,
		}

		dto := mapUserToUserDTO(user)

		assert.Equal(t, user.UserID, dto.UserID)
		assert.Equal(t, user.Username, dto.Username)
		assert.Equal(t, user.Email, dto.Email)
		assert.Equal(t, string(user.Role), dto.Role)
		assert.Equal(t, user.Region, dto.Region)
		assert.Equal(t, user.Active, dto.Active)
		require.NotNil(t, dto.LastAbsenceAt)
		assert.Equal(t, absenceAt, *dto.LastAbsenceAt)
	})

	t.Run("maps a user with last absence at nil", func(t *testing.T) {
		user := domain.User{UserID: "user-1"}

		dto := mapUserToUserDTO(user)

		assert.Nil(t, dto.LastAbsenceAt)
	})
}

func TestMapUserDTOToUser(t *testing.T) {
	t.Run("maps a dto with last absence at set", func(t *testing.T) {
		now := time.Now().UTC()
		absenceAt := now.Add(-time.Hour)
		sessionToken := "token"
		dto := UserDTO{
			UserID:        "user-1",
			Username:      "johndoe",
			FirstName:     "John",
			LastName:      "Doe",
			Email:         "john@doe.com",
			Role:          string(domain.OPERATOR),
			Region:        1,
			SessionToken:  &sessionToken,
			Active:        true,
			LastAbsenceAt: &absenceAt,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		user := mapUserDTOToUser(dto)

		assert.Equal(t, dto.UserID, user.UserID)
		assert.Equal(t, domain.UserRole(dto.Role), user.Role)
		assert.Equal(t, sessionToken, user.SessionToken)
		require.NotNil(t, user.LastAbsenceAt)
		assert.Equal(t, absenceAt, *user.LastAbsenceAt)
	})

	t.Run("maps a dto with last absence at nil", func(t *testing.T) {
		dto := UserDTO{UserID: "user-1"}

		user := mapUserDTOToUser(dto)

		assert.Nil(t, user.LastAbsenceAt)
	})
}

func TestMapUserDTOsToUsers(t *testing.T) {
	dtos := []UserDTO{
		{UserID: "user-1"},
		{UserID: "user-2"},
	}

	users := mapUserDTOsToUsers(dtos)

	assert.Len(t, users, 2)
	assert.Equal(t, "user-1", users[0].UserID)
	assert.Equal(t, "user-2", users[1].UserID)
}
