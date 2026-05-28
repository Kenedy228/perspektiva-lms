package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
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

func TestRequireStudent(t *testing.T) {
	tests := []struct {
		name    string
		actor   role.Role
		wantErr bool
	}{
		{name: "студент", actor: role.NewStudent()},
		{name: "администратор", actor: role.NewAdmin(), wantErr: true},
		{name: "создатель", actor: role.NewCreator(), wantErr: true},
		{name: "организация", actor: role.NewOrganization(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.RequireStudent(tt.actor)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrForbidden)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequireProgressAccess(t *testing.T) {
	tests := []struct {
		name    string
		actor   role.Role
		wantErr bool
	}{
		{name: "администратор", actor: role.NewAdmin()},
		{name: "организация", actor: role.NewOrganization()},
		{name: "студент", actor: role.NewStudent()},
		{name: "создатель", actor: role.NewCreator(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.RequireProgressAccess(tt.actor)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrForbidden)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequireOrganization(t *testing.T) {
	tests := []struct {
		name    string
		actor   role.Role
		wantErr bool
	}{
		{name: "администратор", actor: role.NewAdmin()},
		{name: "организация", actor: role.NewOrganization()},
		{name: "студент", actor: role.NewStudent(), wantErr: true},
		{name: "создатель", actor: role.NewCreator(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.RequireOrganization(tt.actor)
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
	tests := []struct {
		name       string
		limit      int
		offset     int
		wantLimit  int
		wantOffset int
		wantErr    bool
		errTarget  error
	}{
		{name: "корректные значения", limit: 50, offset: 10, wantLimit: 50, wantOffset: 10},
		{name: "limit=0 применяет DefaultLimit", limit: 0, offset: 5, wantLimit: common.DefaultLimit, wantOffset: 5},
		{name: "отрицательный offset", limit: 10, offset: -1, wantErr: true, errTarget: common.ErrInvalidInput},
		{name: "отрицательный limit", limit: -1, offset: 0, wantErr: true, errTarget: common.ErrInvalidInput},
		{name: "превышение MaxLimit", limit: common.MaxLimit + 1, offset: 0, wantErr: true, errTarget: common.ErrInvalidInput},
		{name: "ровно MaxLimit", limit: common.MaxLimit, offset: 0, wantLimit: common.MaxLimit, wantOffset: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l, o, err := common.NormalizePagination(tt.limit, tt.offset)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.errTarget)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantLimit, l)
			assert.Equal(t, tt.wantOffset, o)
		})
	}
}

func TestNormalizeSearchText(t *testing.T) {
	assert.Equal(t, "курс по Go", common.NormalizeSearchText("  курс  по   Go  "))
	assert.Equal(t, "", common.NormalizeSearchText("   "))
	assert.Equal(t, "Go", common.NormalizeSearchText("Go"))
}
