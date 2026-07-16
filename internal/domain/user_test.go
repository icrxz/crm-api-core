package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("creates a user successfully", func(t *testing.T) {
		user, err := NewUser("John", "Doe", "john@doe.com", "s3cr3t", "author-1", "johndoe", OPERATOR, 1)

		require.NoError(t, err)
		assert.NotEmpty(t, user.UserID)
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "john@doe.com", user.Email)
		assert.Equal(t, "johndoe", user.Username)
		assert.Equal(t, OPERATOR, user.Role)
		assert.Equal(t, 1, user.Region)
		assert.NotEqual(t, "s3cr3t", user.Password)
		assert.True(t, user.Active)
		assert.Equal(t, "author-1", user.CreatedBy)
		assert.Equal(t, "author-1", user.UpdatedBy)
		assert.False(t, user.CreatedAt.IsZero())
		assert.Equal(t, user.CreatedAt, user.UpdatedAt)
	})

	t.Run("hashes the password so it can be compared later", func(t *testing.T) {
		user, err := NewUser("John", "Doe", "john@doe.com", "s3cr3t", "author-1", "johndoe", OPERATOR, 1)

		require.NoError(t, err)
		assert.True(t, user.ComparePassword("s3cr3t"))
		assert.False(t, user.ComparePassword("wrong-password"))
	})
}

func TestUser_SetPassword(t *testing.T) {
	user, err := NewUser("John", "Doe", "john@doe.com", "old-password", "author-1", "johndoe", OPERATOR, 1)
	require.NoError(t, err)

	oldHash := user.Password

	err = user.SetPassword("new-password")

	require.NoError(t, err)
	assert.NotEqual(t, oldHash, user.Password)
	assert.True(t, user.ComparePassword("new-password"))
	assert.False(t, user.ComparePassword("old-password"))
}

func TestUser_MergeUpdate(t *testing.T) {
	t.Run("updates only the provided fields", func(t *testing.T) {
		user, err := NewUser("John", "Doe", "john@doe.com", "s3cr3t", "author-1", "johndoe", OPERATOR, 1)
		require.NoError(t, err)

		newFirstName := "Jane"
		newActive := false

		user.MergeUpdate(UserUpdate{
			FirstName: &newFirstName,
			Active:    &newActive,
		}, "author-2")

		assert.Equal(t, newFirstName, user.FirstName)
		assert.False(t, user.Active)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, "author-2", user.UpdatedBy)
	})

	t.Run("updates the remaining fields", func(t *testing.T) {
		user, err := NewUser("John", "Doe", "john@doe.com", "s3cr3t", "author-1", "johndoe", OPERATOR, 1)
		require.NoError(t, err)

		newLastName := "Smith"
		newEmail := "jane@smith.com"
		newLastLoggedIP := "127.0.0.1"
		newSessionToken := "session-token"
		newRegion := 2
		newRole := ADMIN

		user.MergeUpdate(UserUpdate{
			LastName:     &newLastName,
			Email:        &newEmail,
			LastLoggedIP: &newLastLoggedIP,
			SessionToken: &newSessionToken,
			Region:       &newRegion,
			Role:         &newRole,
		}, "author-2")

		assert.Equal(t, newLastName, user.LastName)
		assert.Equal(t, newEmail, user.Email)
		assert.Equal(t, newLastLoggedIP, user.LastLoggedIP)
		assert.Equal(t, newSessionToken, user.SessionToken)
		assert.Equal(t, newRegion, user.Region)
		assert.Equal(t, newRole, user.Role)
	})

	t.Run("leaves fields untouched when update is empty", func(t *testing.T) {
		user, err := NewUser("John", "Doe", "john@doe.com", "s3cr3t", "author-1", "johndoe", OPERATOR, 1)
		require.NoError(t, err)

		oldPassword := user.Password

		user.MergeUpdate(UserUpdate{}, "author-2")

		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, oldPassword, user.Password)
	})

	t.Run("does not change updated by when author is empty", func(t *testing.T) {
		user, err := NewUser("John", "Doe", "john@doe.com", "s3cr3t", "author-1", "johndoe", OPERATOR, 1)
		require.NoError(t, err)

		user.MergeUpdate(UserUpdate{}, "")

		assert.Equal(t, "author-1", user.UpdatedBy)
	})
}
