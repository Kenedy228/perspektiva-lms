package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
)

func TestRequireAdmin(t *testing.T) {
	assert.NoError(t, common.RequireAdmin(role.NewAdmin()))
	assert.ErrorIs(t, common.RequireAdmin(role.NewStudent()), common.ErrForbidden)
}

func TestNormalizePagination(t *testing.T) {
	limit, offset, err := common.NormalizePagination(0, 5)
	assert.NoError(t, err)
	assert.Equal(t, common.DefaultLimit, limit)
	assert.Equal(t, 5, offset)

	_, _, err = common.NormalizePagination(-1, 0)
	assert.ErrorIs(t, err, common.ErrInvalidInput)

	_, _, err = common.NormalizePagination(1, -1)
	assert.ErrorIs(t, err, common.ErrInvalidInput)

	_, _, err = common.NormalizePagination(common.MaxLimit+1, 0)
	assert.ErrorIs(t, err, common.ErrInvalidInput)
}

func TestNormalizeSearchText(t *testing.T) {
	assert.Equal(t, "admin user", common.NormalizeSearchText("  admin   user  "))
}
