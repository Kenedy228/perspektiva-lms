package account

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_validateID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "UUID.Nil",
			args: args{
				id: uuid.Nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "не Nil uuid",
			args: args{
				id: uuid.New(),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateID(tt.args.id), fmt.Sprintf("validateID(%v)", tt.args.id))
		})
	}
}

func Test_validateLogin(t *testing.T) {
	type args struct {
		l login.Login
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "пустой логин",
			args: args{
				l: login.Login{},
			},
			wantErr: assert.Error,
		},
		{
			name: "непустой логин",
			args: args{
				l: loginFixture(t),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateLogin(tt.args.l), fmt.Sprintf("validateLogin(%v)", tt.args.l))
		})
	}
}

func Test_validatePasswordHash(t *testing.T) {
	type args struct {
		h passhash.Hash
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "пустой хеш",
			args: args{
				h: passhash.Hash{},
			},
			wantErr: assert.Error,
		},
		{
			name: "непустой хеш",
			args: args{
				h: passhashFixture(t),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validatePasswordHash(tt.args.h), fmt.Sprintf("validatePasswordHash(%v)", tt.args.h))
		})
	}
}

func Test_validatePersonID(t *testing.T) {
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "uuid.Nil",
			args: args{
				id: uuid.Nil,
			},
			wantErr: assert.Error,
		},
		{
			name: "существующий id",
			args: args{
				id: uuid.New(),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validatePersonID(tt.args.id), fmt.Sprintf("validatePersonID(%v)", tt.args.id))
		})
	}
}

func Test_validateRole(t *testing.T) {
	type args struct {
		r role.Role
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "несуществующая роль",
			args: args{
				r: role.Role{},
			},
			wantErr: assert.Error,
		},
		{
			name: "существующая роль",
			args: args{
				r: role.NewAdmin(),
			},
			wantErr: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateRole(tt.args.r), fmt.Sprintf("validateRole(%v)", tt.args.r))
		})
	}
}

func Test_validateStatus(t *testing.T) {
	type args struct {
		s Status
	}

	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "несуществующий статус",
			args: args{
				s: Status("undefined"),
			},
			wantErr: assert.Error,
		},
		{
			name: "существующий статус",
			args: args{
				s: StatusActive,
			},
			wantErr: assert.NoError,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, validateStatus(tt.args.s), fmt.Sprintf("validateStatus(%v)", tt.args.s))
		})
	}
}
