package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestRequireAdmin(t *testing.T) {
	assert.NoError(t, common.RequireAdmin(role.NewAdmin()))

	err := common.RequireAdmin(role.NewStudent())
	assert.ErrorIs(t, err, common.ErrForbidden)
}

func TestNormalizePagination(t *testing.T) {
	t.Run("defaults zero limit", func(t *testing.T) {
		limit, offset, err := common.NormalizePagination(0, 5)
		assert.NoError(t, err)
		assert.Equal(t, common.DefaultLimit, limit)
		assert.Equal(t, 5, offset)
	})

	t.Run("rejects invalid values", func(t *testing.T) {
		_, _, err := common.NormalizePagination(-1, 0)
		assert.ErrorIs(t, err, common.ErrInvalidInput)

		_, _, err = common.NormalizePagination(1, -1)
		assert.ErrorIs(t, err, common.ErrInvalidInput)

		_, _, err = common.NormalizePagination(common.MaxLimit+1, 0)
		assert.ErrorIs(t, err, common.ErrInvalidInput)
	})
}

func TestNormalizeSearchValues(t *testing.T) {
	assert.Equal(t, "Иванов Иван", common.NormalizeSearchText("  Иванов   Иван  "))
	assert.Equal(t, "11223344595", common.NormalizeSNILSSearch("112-233-445 95"))
}
