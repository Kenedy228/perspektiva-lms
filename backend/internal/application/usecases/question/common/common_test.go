package common_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequireAuthor(t *testing.T) {
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
			err := common.RequireAuthor(tt.actor)
			if tt.wantErr {
				require.Error(t, err)
				assert.ErrorIs(t, err, common.ErrForbidden)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
