package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireManager(t *testing.T) {
	tests := []struct {
		name    string
		actor   role.Role
		wantErr bool
	}{
		{name: "администратор", actor: role.NewAdmin()},
		{name: "создатель", actor: role.NewCreator()},
		{name: "студент", actor: role.NewStudent(), wantErr: true},
		{name: "организация", actor: role.NewOrganization(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.RequireManager(tt.actor)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrForbidden)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizePagination(t *testing.T) {
	t.Run("применяет лимит по умолчанию при limit=0", func(t *testing.T) {
		limit, offset, err := common.NormalizePagination(0, 10)

		require.NoError(t, err)
		assert.Equal(t, common.DefaultLimit, limit)
		assert.Equal(t, 10, offset)
	})

	t.Run("передаёт явно указанный корректный лимит", func(t *testing.T) {
		limit, offset, err := common.NormalizePagination(50, 5)

		require.NoError(t, err)
		assert.Equal(t, 50, limit)
		assert.Equal(t, 5, offset)
	})

	t.Run("отрицательное смещение", func(t *testing.T) {
		_, _, err := common.NormalizePagination(10, -1)

		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("отрицательный лимит", func(t *testing.T) {
		_, _, err := common.NormalizePagination(-1, 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("превышение максимального лимита", func(t *testing.T) {
		_, _, err := common.NormalizePagination(common.MaxLimit+1, 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})
}

func TestNormalizeSearchText(t *testing.T) {
	assert.Equal(t, "банк вопросов", common.NormalizeSearchText("  банк   вопросов  "))
	assert.Equal(t, "", common.NormalizeSearchText("   "))
	assert.Equal(t, "abc", common.NormalizeSearchText("abc"))
}
