package role_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("creates role from valid type", func(t *testing.T) {
		got, err := role.New(role.TypeCreator)

		require.NoError(t, err)
		assert.Equal(t, role.TypeCreator, got.Kind())
		assert.False(t, got.IsZero())
	})

	t.Run("rejects invalid type", func(t *testing.T) {
		got, err := role.New(role.Type("unknown"))

		require.Error(t, err)
		assert.True(t, errors.Is(err, role.ErrInvalid))
		assert.True(t, got.IsZero())
	})
}

func TestConstructors(t *testing.T) {
	tests := []struct {
		name string
		make func() role.Role
		want role.Type
	}{
		{name: "admin", make: role.NewAdmin, want: role.TypeAdmin},
		{name: "creator", make: role.NewCreator, want: role.TypeCreator},
		{name: "student", make: role.NewStudent, want: role.TypeStudent},
		{name: "organization", make: role.NewOrganization, want: role.TypeOrganization},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.make()

			assert.Equal(t, tt.want, got.Kind())
			assert.False(t, got.IsZero())
		})
	}
}

func TestRoleAllows(t *testing.T) {
	creator := role.NewCreator()

	assert.True(t, creator.Allows(role.ResourceCourse, role.ActionPublish))
	assert.False(t, creator.Allows(role.ResourceAuditLog, role.ActionRead))
}
