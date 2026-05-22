//go:build legacy
// +build legacy

package role_test

import (
	"testing"

	role2 "gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestTypeTitle(t *testing.T) {
	t.Run("для фиксированных значений возвращает непустую строку", func(t *testing.T) {
		//Assert
		assert.NotEmpty(t, role2.TypeAdmin.Title())
		assert.NotEmpty(t, role2.TypeOrganization.Title())
		assert.NotEmpty(t, role2.TypeStudent.Title())
		assert.NotEmpty(t, role2.TypeCreator.Title())
	})

	t.Run("для нефиксированных значений возвращает пустую строку", func(t *testing.T) {
		//Assert
		assert.Empty(t, role2.Type("unexisting").Title())
	})
}

func TestTypeAllows(t *testing.T) {
	tc := []struct {
		name    string
		rType   role2.Type
		res     role2.Resource
		action  role2.Action
		allowed bool
	}{
		{
			name:    "администратор разрешает доступ к существующему ресурсу и действию",
			rType:   role2.TypeAdmin,
			res:     role2.ResourceUser,
			action:  role2.ActionWrite,
			allowed: true,
		},
		{
			name:    "администратор запрещает доступ к несуществующему ресурсу с существующим действием",
			rType:   role2.TypeAdmin,
			res:     role2.Resource("unexisting"),
			action:  role2.ActionWrite,
			allowed: false,
		},
		{
			name:    "администратор запрещает несуществующее действие с существующим ресурсом",
			rType:   role2.TypeAdmin,
			res:     role2.ResourceUser,
			action:  role2.Action("unexisting"),
			allowed: false,
		},
		{
			name:    "несуществующий тип запрещает доступ к существующему ресурс и действию",
			rType:   role2.Type("incorrect"),
			res:     role2.ResourceUser,
			action:  role2.ActionRead,
			allowed: false,
		},
		{
			name:    "студент запрещает доступ к ресурсу, который не привязан к нему",
			rType:   role2.TypeStudent,
			res:     role2.ResourceUser,
			action:  role2.ActionWrite,
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
