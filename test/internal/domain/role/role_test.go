package role_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestAllowsRole(t *testing.T) {
	t.Run("existing resource, existing action, correct type", func(t *testing.T) {
		//Arrange
		r := role.NewAdmin()
		res := role.ResourceUser
		action := role.ActionWrite

		//Assert
		assert.True(t, r.Allows(res, action))
	})

	t.Run("non-existing resource, existing action, correct type", func(t *testing.T) {
		//Arrange
		r := role.NewAdmin()
		res := role.Resource("non-existing")
		action := role.ActionWrite

		//Assert
		assert.False(t, r.Allows(res, action))
	})

	t.Run("existing resource, non-existing action, correct type", func(t *testing.T) {
		//Arrange
		r := role.NewAdmin()
		res := role.Resource("non-existing")
		action := role.ActionWrite

		//Assert
		assert.False(t, r.Allows(res, action))
	})

	t.Run("existing resource, non-existing action, correct type", func(t *testing.T) {
		//Arrange
		r := role.NewAdmin()
		res := role.ResourceUser
		action := role.Action("non-existing")

		//Assert
		assert.False(t, r.Allows(res, action))
	})

	t.Run("existing resource, existing action, type without resource", func(t *testing.T) {
		//Arrange
		r := role.NewStudent()
		res := role.ResourceUser
		action := role.ActionWrite

		//Assert
		assert.False(t, r.Allows(res, action))
	})
}
