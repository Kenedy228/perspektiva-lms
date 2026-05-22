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

func TestNew(t *testing.T) {
	// Arrange
	org, err := organization.New(innFixture(), nameFixture())
	require.NoError(t, err)

	// Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, org.ID())
	assert.Equal(t, innFixture(), org.INN())
	assert.Equal(t, nameFixture(), org.Name())
}

func TestChangeINN(t *testing.T) {
	// Arrange
	org, err := organization.New(innFixture(), nameFixture())
	require.NoError(t, err)

	// Act
	i, err := inn2.New("7728168971", inn2.TypeOrganization)
	require.NoError(t, err)
	org.ChangeINN(i)

	// Assert
	assert.Equal(t, i, org.INN())
}

func TestRename(t *testing.T) {
	// Arrange
	org, err := organization.New(innFixture(), nameFixture())
	require.NoError(t, err)

	// Act
	n, err := name.New("ООО 'Спартак'")
	require.NoError(t, err)

	org.Rename(n)

	// Assert
	assert.Equal(t, n, org.Name())
}

// Fixtures

func innFixture() inn2.INN {
	i, _ := inn2.New("1030000000", inn2.TypeOrganization)
	return i
}

func nameFixture() name.Name {
	n, _ := name.New("ООО 'Ромашка'")
	return n
}
