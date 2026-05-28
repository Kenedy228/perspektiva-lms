package role_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseType(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    role.Type
		wantErr bool
	}{
		{
			name:  "parses valid role type",
			value: "student",
			want:  role.TypeStudent,
		},
		{
			name:  "trims spaces around valid role type",
			value: " admin ",
			want:  role.TypeAdmin,
		},
		{
			name:    "rejects unknown role type",
			value:   "moderator",
			wantErr: true,
		},
		{
			name:    "rejects empty role type",
			value:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := role.ParseType(tt.value)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, role.ErrInvalid))
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTypeAllows(t *testing.T) {
	tests := []struct {
		name     string
		roleType role.Type
		resource role.Resource
		action   role.Action
		want     bool
	}{
		{
			name:     "admin can read audit logs",
			roleType: role.TypeAdmin,
			resource: role.ResourceAuditLog,
			action:   role.ActionRead,
			want:     true,
		},
		{
			name:     "creator can grade submissions",
			roleType: role.TypeCreator,
			resource: role.ResourceSubmission,
			action:   role.ActionGrade,
			want:     true,
		},
		{
			name:     "student can submit an attempt",
			roleType: role.TypeStudent,
			resource: role.ResourceAttempt,
			action:   role.ActionSubmit,
			want:     true,
		},
		{
			name:     "student cannot grade",
			roleType: role.TypeStudent,
			resource: role.ResourceGrade,
			action:   role.ActionGrade,
			want:     false,
		},
		{
			name:     "organization can read certificates",
			roleType: role.TypeOrganization,
			resource: role.ResourceCertificate,
			action:   role.ActionRead,
			want:     true,
		},
		{
			name:     "creator can create blocks",
			roleType: role.TypeCreator,
			resource: role.ResourceBlock,
			action:   role.ActionCreate,
			want:     true,
		},
		{
			name:     "creator can create elements",
			roleType: role.TypeCreator,
			resource: role.ResourceElement,
			action:   role.ActionCreate,
			want:     true,
		},
		{
			name:     "creator can manage bank",
			roleType: role.TypeCreator,
			resource: role.ResourceBank,
			action:   role.ActionDelete,
			want:     true,
		},
		{
			name:     "student can read blocks",
			roleType: role.TypeStudent,
			resource: role.ResourceBlock,
			action:   role.ActionRead,
			want:     true,
		},
		{
			name:     "student cannot write blocks",
			roleType: role.TypeStudent,
			resource: role.ResourceBlock,
			action:   role.ActionCreate,
			want:     false,
		},
		{
			name:     "student cannot access banks",
			roleType: role.TypeStudent,
			resource: role.ResourceBank,
			action:   role.ActionRead,
			want:     false,
		},
		{
			name:     "organization cannot access blocks",
			roleType: role.TypeOrganization,
			resource: role.ResourceBlock,
			action:   role.ActionRead,
			want:     false,
		},
		{
			name:     "unknown role denies access",
			roleType: role.Type("unknown"),
			resource: role.ResourceCourse,
			action:   role.ActionRead,
			want:     false,
		},
		{
			name:     "unknown resource denies access",
			roleType: role.TypeAdmin,
			resource: role.Resource("unknown"),
			action:   role.ActionRead,
			want:     false,
		},
		{
			name:     "unknown action denies access",
			roleType: role.TypeAdmin,
			resource: role.ResourceCourse,
			action:   role.Action("unknown"),
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.roleType.Allows(tt.resource, tt.action))
		})
	}
}

func TestTypeTitle(t *testing.T) {
	tests := []role.Type{
		role.TypeAdmin,
		role.TypeCreator,
		role.TypeStudent,
		role.TypeOrganization,
	}

	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			assert.NotEmpty(t, tt.Title())
		})
	}

	assert.Empty(t, role.Type("unknown").Title())
}
