package account_test

import (
	"gitflic.ru/lms/internal/domain/account"
	"gitflic.ru/lms/internal/domain/account/login"
	"gitflic.ru/lms/internal/domain/account/passhash"
	"gitflic.ru/lms/internal/domain/role"
	"github.com/google/uuid"
)

type accountBuilder struct {
	login        login.Login
	passwordHash passhash.Hash
	role         role.Role
	personID     uuid.UUID
}

func newAccountBuilder() *accountBuilder {
	return &accountBuilder{}
}

func (b *accountBuilder) withLogin() *accountBuilder {
	b.login = loginFixture()
	return b
}

func (b *accountBuilder) withPasswordHash() *accountBuilder {
	b.passwordHash = hashFixture()
	return b
}

func (b *accountBuilder) withRole() *accountBuilder {
	b.role = roleFixture()
	return b
}

func (b *accountBuilder) withPersonID(id uuid.UUID) *accountBuilder {
	b.personID = id
	return b
}

func (b *accountBuilder) build() (*account.Account, error) {
	params := account.Params{
		Login:        b.login,
		PasswordHash: b.passwordHash,
		Role:         b.role,
		PersonID:     b.personID,
	}

	return account.New(params)
}
