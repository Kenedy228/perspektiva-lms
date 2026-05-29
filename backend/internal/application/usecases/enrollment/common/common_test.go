package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/enrollment/common"
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
		{name: "создатель", actor: role.NewCreator(), wantErr: true},
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

func TestRequireViewer(t *testing.T) {
	tests := []struct {
		name    string
		actor   role.Role
		wantErr bool
	}{
		{name: "администратор", actor: role.NewAdmin()},
		{name: "студент", actor: role.NewStudent()},
		{name: "создатель", actor: role.NewCreator(), wantErr: true},
		{name: "организация", actor: role.NewOrganization(), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := common.RequireViewer(tt.actor)
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
