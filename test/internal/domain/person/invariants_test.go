//go:build legacy
// +build legacy

package person_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/person"
	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRejectsEmptyName(t *testing.T) {
	p, err := person.New(name.Name{})

	require.Error(t, err)
	assert.ErrorIs(t, err, person.ErrInvalid)
	assert.Nil(t, p)
}

func TestRestorePerson(t *testing.T) {
	t.Run("restores person with existing id and profile", func(t *testing.T) {
		id := uuid.New()
		prof := profileFixture()

		p, err := person.Restore(id, nameFixture(), &prof)

		require.NoError(t, err)
		assert.Equal(t, id, p.ID())
		assert.True(t, p.HasProfile())
	})

	t.Run("rejects empty id", func(t *testing.T) {
		p, err := person.Restore(uuid.Nil, nameFixture(), nil)

		require.Error(t, err)
		assert.ErrorIs(t, err, person.ErrInvalid)
		assert.Nil(t, p)
	})
}

func TestAttachProfileRejectsEmptyProfile(t *testing.T) {
	p, err := person.New(nameFixture())
	require.NoError(t, err)

	err = p.AttachProfile(profile.Profile{})

	require.Error(t, err)
	assert.ErrorIs(t, err, person.ErrInvalid)
	assert.False(t, p.HasProfile())
}

func TestProfileNewRejectsMissingRequiredValues(t *testing.T) {
	prof, err := profile.New(profile.Params{})

	require.Error(t, err)
	assert.ErrorIs(t, err, profile.ErrInvalid)
	assert.True(t, prof.IsZero())
}
