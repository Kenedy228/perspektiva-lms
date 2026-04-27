package account_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/account"
	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type accountBuilder struct {
	login        string
	passwordHash string
	role         role.Role
	personID     uuid.UUID
}

func newAccountBuilder() *accountBuilder {
	return &accountBuilder{}
}

func (b *accountBuilder) withLogin(login string) *accountBuilder {
	b.login = login
	return b
}

func (b *accountBuilder) withPasswordHash(hash string) *accountBuilder {
	b.passwordHash = hash
	return b
}

func (b *accountBuilder) withRole(r role.Role) *accountBuilder {
	b.role = r
	return b
}

func (b *accountBuilder) withPersonID(id uuid.UUID) *accountBuilder {
	b.personID = id
	return b
}

func (b *accountBuilder) build(t *testing.T, wantErr error) *account.Account {
	t.Helper()

	params := account.Params{
		Login:        b.login,
		PasswordHash: b.passwordHash,
		Role:         b.role,
		PersonID:     b.personID,
	}

	acc, err := account.New(params)
	assert.ErrorIs(t, err, wantErr)

	return acc
}
