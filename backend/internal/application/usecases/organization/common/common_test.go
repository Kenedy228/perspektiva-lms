package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAdmin(t *testing.T) {
	require.NoError(t, common.RequireAdmin(role.NewAdmin()))

	err := common.RequireAdmin(role.NewStudent())
	require.Error(t, err)
	assert.ErrorIs(t, err, common.ErrForbidden)
}

func TestNormalizePagination(t *testing.T) {
	t.Run("applies default limit", func(t *testing.T) {
		limit, offset, err := common.NormalizePagination(0, 7)

		require.NoError(t, err)
		assert.Equal(t, common.DefaultLimit, limit)
		assert.Equal(t, 7, offset)
	})

	t.Run("rejects negative offset", func(t *testing.T) {
		_, _, err := common.NormalizePagination(10, -1)

		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})

	t.Run("rejects too large limit", func(t *testing.T) {
		_, _, err := common.NormalizePagination(common.MaxLimit+1, 0)

		require.Error(t, err)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})
}
