package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/attempt/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
