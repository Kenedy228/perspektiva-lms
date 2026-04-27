package role_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("должен создать через специальный конструктор роль с соответствующим типом", func(t *testing.T) {
		//Arrange
		admin := role.NewAdmin()

		//Assert
		assert.Equal(t, admin.Kind(), role.TypeAdmin)
	})
}

func TestRoleAllows_Admin(t *testing.T) {
	tc := []struct {
		name    string
		rType   role.Type
		res     role.Resource
		action  role.Action
		allowed bool
	}{
		{
			name:    "администратор разрешает доступ к существующему ресурсу и действию",
			res:     role.ResourceUser,
			action:  role.ActionWrite,
			allowed: true,
		},
		{
			name:    "администратор запрещает доступ к несуществующему ресурсу с существующим действием",
			res:     role.Resource("unexisting"),
			action:  role.ActionWrite,
			allowed: false,
		},
		{
			name:    "администратор запрещает несуществующее действие с существующим ресурсом",
			res:     role.ResourceUser,
			action:  role.Action("unexisting"),
			allowed: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			admin := role.NewAdmin()

			//Act
			allowed := admin.Allows(tt.res, tt.action)

			//Assert
			assert.Equal(t, allowed, tt.allowed)
		})
	}
}
