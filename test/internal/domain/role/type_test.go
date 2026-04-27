package role_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	t.Run("for existing types should return non-empty title", func(t *testing.T) {
		//Arrange
		admin := role.TypeAdmin
		organization := role.TypeOrganization
		student := role.TypeStudent
		creator := role.TypeCreator

		//Assert
		assert.NotEmpty(t, admin.Title())
		assert.NotEmpty(t, organization.Title())
		assert.NotEmpty(t, student.Title())
		assert.NotEmpty(t, creator.Title())
	})

	t.Run("for non-existing types should return empty", func(t *testing.T) {
		//Arrange
		nonExisting := role.Type("non-existing")

		//Assert
		assert.Empty(t, nonExisting.Title())
	})
}

func TestAllows(t *testing.T) {
	t.Run("existing resource, existing action, correct type", func(t *testing.T) {
		//Arrange
		rType := role.TypeAdmin
		res := role.ResourceUser
		action := role.ActionWrite

		//Assert
		assert.True(t, rType.Allows(res, action))
	})

	t.Run("non-existing resource, existing action, correct type", func(t *testing.T) {
		//Arrange
		rType := role.TypeAdmin
		res := role.Resource("non-existing")
		action := role.ActionWrite

		//Assert
		assert.False(t, rType.Allows(res, action))
	})

	t.Run("existing resource, non-existing action, correct type", func(t *testing.T) {
		//Arrange
		rType := role.TypeAdmin
		res := role.Resource("non-existing")
		action := role.ActionWrite

		//Assert
		assert.False(t, rType.Allows(res, action))
	})

	t.Run("existing resource, non-existing action, correct type", func(t *testing.T) {
		//Arrange
		rType := role.TypeAdmin
		res := role.ResourceUser
		action := role.Action("non-existing")

		//Assert
		assert.False(t, rType.Allows(res, action))
	})

	t.Run("existing resource, existing action, incorrect type", func(t *testing.T) {
		//Arrange
		rType := role.Type("incorrect")
		res := role.ResourceUser
		action := role.ActionRead

		//Assert
		assert.False(t, rType.Allows(res, action))
	})

	t.Run("existing resource, existing action, type without resource", func(t *testing.T) {
		//Arrange
		rType := role.TypeStudent
		res := role.ResourceUser
		action := role.ActionWrite

		//Assert
		assert.False(t, rType.Allows(res, action))
	})
}
