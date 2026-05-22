package account

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func Test_ensureNotDeleted(t *testing.T) {
	type args struct {
		a *Account
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "аккаунт не удален",
			args: args{
				a: &Account{
					id:           uuid.UUID{},
					login:        login.Login{},
					passwordHash: passhash.Hash{},
					role:         role.Role{},
					personID:     uuid.UUID{},
					status:       StatusActive,
				},
			},
			wantErr: false,
		},
		{
			name: "аккаунт удален",
			args: args{
				a: &Account{
					id:           uuid.UUID{},
					login:        login.Login{},
					passwordHash: passhash.Hash{},
					role:         role.Role{},
					personID:     uuid.UUID{},
					status:       StatusDeleted,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ensureNotDeleted(tt.args.a); (err != nil) != tt.wantErr {
				t.Errorf("ensureNotDeleted() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
