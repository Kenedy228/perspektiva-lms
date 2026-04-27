package role_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestTypeTitle(t *testing.T) {
	t.Run("для фиксированных значений возвращает непустую строку", func(t *testing.T) {
		//Assert
		assert.NotEmpty(t, role.TypeAdmin.Title())
		assert.NotEmpty(t, role.TypeOrganization.Title())
		assert.NotEmpty(t, role.TypeStudent.Title())
		assert.NotEmpty(t, role.TypeCreator.Title())
	})

	t.Run("для нефиксированных значений возвращает пустую строку", func(t *testing.T) {
		//Assert
		assert.Empty(t, role.Type("unexisting").Title())
	})
}

func TestTypeAllows(t *testing.T) {
	tc := []struct {
		name    string
		rType   role.Type
		res     role.Resource
		action  role.Action
		allowed bool
	}{
		{
			name:    "администратор разрешает доступ к существующему ресурсу и действию",
			rType:   role.TypeAdmin,
			res:     role.ResourceUser,
			action:  role.ActionWrite,
			allowed: true,
		},
		{
			name:    "администратор запрещает доступ к несуществующему ресурсу с существующим действием",
			rType:   role.TypeAdmin,
			res:     role.Resource("unexisting"),
			action:  role.ActionWrite,
			allowed: false,
		},
		{
			name:    "администратор запрещает несуществующее действие с существующим ресурсом",
			rType:   role.TypeAdmin,
			res:     role.ResourceUser,
			action:  role.Action("unexisting"),
			allowed: false,
		},
		{
			name:    "несуществующий тип запрещает доступ к существующему ресурс и действию",
			rType:   role.Type("incorrect"),
			res:     role.ResourceUser,
			action:  role.ActionRead,
			allowed: false,
		},
		{
			name:    "студент запрещает доступ к ресурсу, который не привязан к нему",
			rType:   role.TypeStudent,
			res:     role.ResourceUser,
			action:  role.ActionWrite,
			allowed: false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Act
			allowed := tt.rType.Allows(tt.res, tt.action)

			//Assert
			assert.Equal(t, allowed, tt.allowed)
		})
	}
}
