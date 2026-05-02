package organization_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/organization/inn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	org, err := newOrganizationBuilder().
		withInn("1030000000", inn.TypeOrganization).
		withName("name").
		build()

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, org.INN().Code(), "1030000000")
	assert.Equal(t, org.Name().Value(), "name")
}

func TestChangeINN(t *testing.T) {
	//Arrange
	org, err := newOrganizationBuilder().
		withInn("1030000000", inn.TypeOrganization).
		withName("name").
		build()
	require.NoError(t, err)

	//Act
	org.ChangeINN(makeINN("7728168971", inn.TypeOrganization))

	//Assert
	assert.Equal(t, org.INN().Code(), "7728168971")
}

func TestRename(t *testing.T) {
	//Arrange
	org, err := newOrganizationBuilder().
		withInn("1030000000", inn.TypeOrganization).
		withName("name").
		build()
	require.NoError(t, err)

	//Act
	org.Rename(makeName("new name"))

	//Assert
	assert.Equal(t, org.Name().Value(), "new name")
}
