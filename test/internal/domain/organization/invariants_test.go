//go:build legacy
// +build legacy

package organization_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/organization"
	inn2 "gitflic.ru/lms/backend/internal/domain/organization/inn"
	"gitflic.ru/lms/backend/internal/domain/organization/name"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRejectsEmptyValues(t *testing.T) {
	t.Run("rejects empty inn", func(t *testing.T) {
		org, err := organization.New(inn2.INN{}, nameFixture())

		require.Error(t, err)
		assert.ErrorIs(t, err, organization.ErrInvalid)
		assert.Nil(t, org)
	})

	t.Run("rejects empty name", func(t *testing.T) {
		org, err := organization.New(innFixture(), name.Name{})

		require.Error(t, err)
		assert.ErrorIs(t, err, organization.ErrInvalid)
		assert.Nil(t, org)
	})
}

func TestRestoreOrganization(t *testing.T) {
	t.Run("restores organization with existing id", func(t *testing.T) {
		id := uuid.New()

		org, err := organization.Restore(id, innFixture(), nameFixture())

		require.NoError(t, err)
		assert.Equal(t, id, org.ID())
	})

	t.Run("rejects empty id", func(t *testing.T) {
		org, err := organization.Restore(uuid.Nil, innFixture(), nameFixture())

		require.Error(t, err)
		assert.ErrorIs(t, err, organization.ErrInvalid)
		assert.Nil(t, org)
	})
}

func TestMutatorsRejectEmptyValues(t *testing.T) {
	org, err := organization.New(innFixture(), nameFixture())
	require.NoError(t, err)

	prevINN := org.INN()
	prevName := org.Name()

	err = org.ChangeINN(inn2.INN{})
	require.Error(t, err)
	assert.ErrorIs(t, err, organization.ErrInvalid)
	assert.Equal(t, prevINN, org.INN())

	err = org.Rename(name.Name{})
	require.Error(t, err)
	assert.ErrorIs(t, err, organization.ErrInvalid)
	assert.Equal(t, prevName, org.Name())
}
