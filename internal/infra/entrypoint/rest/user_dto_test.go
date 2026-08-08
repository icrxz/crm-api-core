package rest

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNullableTime_UnmarshalJSON(t *testing.T) {
	t.Run("sets value when a valid time is sent", func(t *testing.T) {
		var n NullableTime
		err := json.Unmarshal([]byte(`"2026-08-07T10:00:00Z"`), &n)

		require.NoError(t, err)
		assert.True(t, n.Set)
		require.NotNil(t, n.Value)
		assert.Equal(t, time.Date(2026, 8, 7, 10, 0, 0, 0, time.UTC), *n.Value)
	})

	t.Run("marks set with nil value when sent as null", func(t *testing.T) {
		var n NullableTime
		err := json.Unmarshal([]byte(`null`), &n)

		require.NoError(t, err)
		assert.True(t, n.Set)
		assert.Nil(t, n.Value)
	})

	t.Run("returns error when value is not a valid time", func(t *testing.T) {
		var n NullableTime
		err := json.Unmarshal([]byte(`"not-a-time"`), &n)

		assert.Error(t, err)
	})

	t.Run("leaves unset when the field is absent from the payload", func(t *testing.T) {
		type wrapper struct {
			LastAbsenceAt NullableTime `json:"last_absence_at"`
		}
		var w wrapper
		err := json.Unmarshal([]byte(`{}`), &w)

		require.NoError(t, err)
		assert.False(t, w.LastAbsenceAt.Set)
		assert.Nil(t, w.LastAbsenceAt.Value)
	})
}

func TestNullableTime_MarshalJSON(t *testing.T) {
	t.Run("marshals to null when value is nil", func(t *testing.T) {
		n := NullableTime{Set: true, Value: nil}
		data, err := json.Marshal(n)

		require.NoError(t, err)
		assert.Equal(t, "null", string(data))
	})

	t.Run("marshals the wrapped time when value is set", func(t *testing.T) {
		absenceAt := time.Date(2026, 8, 7, 10, 0, 0, 0, time.UTC)
		n := NullableTime{Set: true, Value: &absenceAt}
		data, err := json.Marshal(n)

		require.NoError(t, err)
		assert.Equal(t, `"2026-08-07T10:00:00Z"`, string(data))
	})
}

func TestMapUpdateUserDTOToUserUpdate(t *testing.T) {
	t.Run("maps last absence at as present with value", func(t *testing.T) {
		absenceAt := time.Date(2026, 8, 7, 10, 0, 0, 0, time.UTC)
		dto := UpdateUserDTO{
			LastAbsenceAt: NullableTime{Set: true, Value: &absenceAt},
		}

		update := mapUpdateUserDTOToUserUpdate(dto)

		assert.Equal(t, domain.OptionalTime{Present: true, Value: &absenceAt}, update.LastAbsenceAt)
	})

	t.Run("maps last absence at as present with nil value when clearing", func(t *testing.T) {
		dto := UpdateUserDTO{
			LastAbsenceAt: NullableTime{Set: true, Value: nil},
		}

		update := mapUpdateUserDTOToUserUpdate(dto)

		assert.Equal(t, domain.OptionalTime{Present: true, Value: nil}, update.LastAbsenceAt)
	})

	t.Run("maps last absence at as not present when field is absent", func(t *testing.T) {
		dto := UpdateUserDTO{}

		update := mapUpdateUserDTOToUserUpdate(dto)

		assert.Equal(t, domain.OptionalTime{Present: false, Value: nil}, update.LastAbsenceAt)
	})
}

func TestMapUserToUserDTO_LastAbsenceAt(t *testing.T) {
	absenceAt := time.Date(2026, 8, 7, 10, 0, 0, 0, time.UTC)
	user := domain.User{UserID: "user-1", LastAbsenceAt: &absenceAt}

	dto := mapUserToUserDTO(user)

	require.NotNil(t, dto.LastAbsenceAt)
	assert.Equal(t, absenceAt, *dto.LastAbsenceAt)
}
