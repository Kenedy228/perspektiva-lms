//go:build legacy
// +build legacy

package role_test

import (
	"testing"

	role2 "gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("должен создать через специальный конструктор роль с соответствующим типом", func(t *testing.T) {
		//Arrange
		admin := role2.NewAdmin()

		//Assert
		assert.Equal(t, admin.Kind(), role2.TypeAdmin)
	})
}

func TestRoleAllows_Admin(t *testing.T) {
	tc := []struct {
		name    string
		rType   role2.Type
		res     role2.Resource
		action  role2.Action
		allowed bool
	}{
		{
			name:    "администратор разрешает доступ к существующему ресурсу и действию",
			res:     role2.ResourceUser,
			action:  role2.ActionWrite,
			allowed: true,
		},
		{
			name:    "администратор запрещает доступ к несуществующему ресурсу с существующим действием",
			res:     role2.Resource("unexisting"),
			action:  role2.ActionWrite,
			allowed: false,
		},
		{
			name:    "администратор запрещает несуществующее действие с существующим ресурсом",
			res:     role2.ResourceUser,
			action:  role2.Action("unexisting"),
			allowed: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			admin := role2.NewAdmin()

			//Act
			allowed := admin.Allows(tt.res, tt.action)

			//Assert
			assert.Equal(t, allowed, tt.allowed)
		})
	}
}
