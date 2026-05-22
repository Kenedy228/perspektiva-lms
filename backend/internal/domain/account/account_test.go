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

func TestAccount_ChangeLogin(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	type args struct {
		login login.Login
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    login.Login
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "логин существует и аккаунт не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			args: args{
				login: loginFixture(t),
			},
			want:    loginFixture(t),
			wantErr: assert.NoError,
		},
		{
			name: "логин не существует и аккаунт не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			args: args{
				login: login.Login{},
			},
			want:    login.Login{},
			wantErr: assert.Error,
		},
		{
			name: "логин существует, но аккаунт удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			args: args{
				login: loginFixture(t),
			},
			want:    login.Login{},
			wantErr: assert.Error,
		},
		{
			name: "логин не существует, но аккаунт удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			args: args{
				login: login.Login{},
			},
			want:    login.Login{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			tt.wantErr(t, a.ChangeLogin(tt.args.login), fmt.Sprintf("ChangeLogin(%v)", tt.args.login))
			assert.Equal(t, tt.want, a.login)
		})
	}
}

func TestAccount_ChangePasswordHash(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	type args struct {
		hash passhash.Hash
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    passhash.Hash
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "хеш существует и аккаунт не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			args: args{
				hash: passhashFixture(t),
			},
			want:    passhashFixture(t),
			wantErr: assert.NoError,
		},
		{
			name: "хеш не существует и аккаунт не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			args: args{
				hash: passhash.Hash{},
			},
			want:    passhash.Hash{},
			wantErr: assert.Error,
		},
		{
			name: "хеш существует, но аккаунт удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			args: args{
				hash: passhashFixture(t),
			},
			want:    passhash.Hash{},
			wantErr: assert.Error,
		},
		{
			name: "хеш не существует, но аккаунт удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			args: args{
				hash: passhash.Hash{},
			},
			want:    passhash.Hash{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			tt.wantErr(t, a.ChangePasswordHash(tt.args.hash), fmt.Sprintf("ChangePasswordHash(%v)", tt.args.hash))
			assert.Equal(t, tt.want, a.passwordHash)
		})
	}
}

func TestAccount_ChangeRole(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	type args struct {
		role role.Role
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    role.Role
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "новая роль существует и аккаунт не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			args: args{
				role: role.NewAdmin(),
			},
			want:    role.NewAdmin(),
			wantErr: assert.NoError,
		},
		{
			name: "роли не существует и аккаунт не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			args: args{
				role: role.Role{},
			},
			want:    role.Role{},
			wantErr: assert.Error,
		},
		{
			name: "роль существует, но аккаунт удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			args: args{
				role.NewAdmin(),
			},
			want:    role.Role{},
			wantErr: assert.Error,
		},
		{
			name: "роль не существует, но аккаунт удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			args: args{
				role.Role{},
			},
			want:    role.Role{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			tt.wantErr(t, a.ChangeRole(tt.args.role), fmt.Sprintf("ChangeRole(%v)", tt.args.role))
			assert.Equal(t, tt.want, a.role)
		})
	}
}

func TestAccount_Delete(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			wantErr: assert.NoError,
		},
		{
			name: "удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			wantErr: assert.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			tt.wantErr(t, a.Delete(), fmt.Sprintf("Delete()"))
			assert.Equal(t, StatusDeleted, a.status)
		})
	}
}

func TestAccount_ID(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "возвращает id как есть",
			fields: fields{
				id:           idFixture,
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       "",
			},
			want: idFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.ID(), "ID()")
		})
	}
}

func TestAccount_IsActive(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "активен",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			want: true,
		},
		{
			name: "неактивен",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.IsActive(), "IsActive()")
		})
	}
}

func TestAccount_IsDeleted(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusDeleted,
			},
			want: true,
		},
		{
			name: "не удален",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.IsDeleted(), "IsDeleted()")
		})
	}
}

func TestAccount_Login(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   login.Login
	}{
		{
			name: "возвращает логин как есть",
			fields: fields{
				id:           uuid.UUID{},
				login:        loginFixture(t),
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       "",
			},
			want: loginFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.Login(), "Login()")
		})
	}
}

func TestAccount_PasswordHash(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   passhash.Hash
	}{
		{
			name: "возвращает хеш пароля как есть",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhashFixture(t),
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       "",
			},
			want: passhashFixture(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.PasswordHash(), "PasswordHash()")
		})
	}
}

func TestAccount_PersonID(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}

	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "возвращает идентификатор обладателя как есть",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     idFixture,
				status:       "",
			},
			want: idFixture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.PersonID(), "PersonID()")
		})
	}
}

func TestAccount_Role(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   role.Role
	}{
		{
			name: "возвращает роль как есть",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.NewAdmin(),
				personID:     uuid.UUID{},
				status:       "",
			},
			want: role.NewAdmin(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.Role(), "Role()")
		})
	}
}

func TestAccount_Status(t *testing.T) {
	type fields struct {
		id           uuid.UUID
		login        login.Login
		passwordHash passhash.Hash
		role         role.Role
		personID     uuid.UUID
		status       Status
	}
	tests := []struct {
		name   string
		fields fields
		want   Status
	}{
		{
			name: "возвращает статус как есть",
			fields: fields{
				id:           uuid.UUID{},
				login:        login.Login{},
				passwordHash: passhash.Hash{},
				role:         role.Role{},
				personID:     uuid.UUID{},
				status:       StatusActive,
			},
			want: StatusActive,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Account{
				id:           tt.fields.id,
				login:        tt.fields.login,
				passwordHash: tt.fields.passwordHash,
				role:         tt.fields.role,
				personID:     tt.fields.personID,
				status:       tt.fields.status,
			}
			assert.Equalf(t, tt.want, a.Status(), "Status()")
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		l        login.Login
		h        passhash.Hash
		r        role.Role
		personID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "без логина",
			args: args{
				l:        login.Login{},
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: idFixture,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без пароля",
			args: args{
				l:        loginFixture(t),
				h:        passhash.Hash{},
				r:        role.NewAdmin(),
				personID: idFixture,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без роли",
			args: args{
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.Role{},
				personID: idFixture,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без идентификатора обладателя",
			args: args{
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: uuid.Nil,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "все поля",
			args: args{
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: idFixture,
			},
			want: &Account{
				id:           uuid.New(),
				login:        loginFixture(t),
				passwordHash: passhashFixture(t),
				role:         role.NewAdmin(),
				personID:     idFixture,
				status:       StatusActive,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.l, tt.args.h, tt.args.r, tt.args.personID)

			if !tt.wantErr(t, err, fmt.Sprintf("New(%v, %v, %v, %v)", tt.args.l, tt.args.h, tt.args.r, tt.args.personID)) {
				return
			}

			if tt.want != nil {
				assert.NotEqual(t, tt.want.id, uuid.Nil)
				assert.Equal(t, tt.want.login, got.login)
				assert.Equal(t, tt.want.passwordHash, got.passwordHash)
				assert.Equal(t, tt.want.role, got.role)
				assert.Equal(t, tt.want.personID, got.personID)
				assert.Equal(t, tt.want.status, got.status)
			}
		})
	}
}

func TestRestore(t *testing.T) {
	type args struct {
		id       uuid.UUID
		l        login.Login
		h        passhash.Hash
		r        role.Role
		personID uuid.UUID
		status   Status
	}
	tests := []struct {
		name    string
		args    args
		want    *Account
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "без идентификатора",
			args: args{
				id:       uuid.Nil,
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: idFixture,
				status:   StatusActive,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без логина",
			args: args{
				id:       idFixture,
				l:        login.Login{},
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: idFixture,
				status:   StatusActive,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без пароля",
			args: args{
				id:       idFixture,
				l:        loginFixture(t),
				h:        passhash.Hash{},
				r:        role.NewAdmin(),
				personID: idFixture,
				status:   StatusActive,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без роли",
			args: args{
				id:       idFixture,
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.Role{},
				personID: idFixture,
				status:   StatusActive,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без идентификатора обладателя",
			args: args{
				id:       idFixture,
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: uuid.Nil,
				status:   StatusActive,
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "без статуса",
			args: args{
				id:       idFixture,
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: idFixture,
				status:   "",
			},
			want:    nil,
			wantErr: assert.Error,
		},
		{
			name: "все поля",
			args: args{
				id:       idFixture,
				l:        loginFixture(t),
				h:        passhashFixture(t),
				r:        role.NewAdmin(),
				personID: idFixture,
				status:   StatusActive,
			},
			want: &Account{
				id:           idFixture,
				login:        loginFixture(t),
				passwordHash: passhashFixture(t),
				role:         role.NewAdmin(),
				personID:     idFixture,
				status:       StatusActive,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.args.id, tt.args.l, tt.args.h, tt.args.r, tt.args.personID, tt.args.status)
			if !tt.wantErr(t, err, fmt.Sprintf("Restore(%v, %v, %v, %v, %v, %v)", tt.args.id, tt.args.l, tt.args.h, tt.args.r, tt.args.personID, tt.args.status)) {
				return
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.id, got.id)
				assert.Equal(t, tt.want.login, got.login)
				assert.Equal(t, tt.want.passwordHash, got.passwordHash)
				assert.Equal(t, tt.want.role, got.role)
				assert.Equal(t, tt.want.personID, got.personID)
				assert.Equal(t, tt.want.status, got.status)
			}
		})
	}
}
