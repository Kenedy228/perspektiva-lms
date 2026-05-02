package account

import (
	"gitflic.ru/lms/internal/domain/account/login"
	"gitflic.ru/lms/internal/domain/account/passhash"
	"gitflic.ru/lms/internal/domain/role"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Account struct {
	id           uuid.UUID
	login        login.Login
	passwordHash passhash.Hash
	role         role.Role
	personID     uuid.UUID
}

func New(params Params) (*Account, error) {
	if err := validatePersonID(params.PersonID); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Account{
		id:           id,
		login:        params.Login,
		passwordHash: params.PasswordHash,
		role:         params.Role,
		personID:     params.PersonID,
	}, nil
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Login() login.Login {
	return a.login
}

func (a *Account) PasswordHash() passhash.Hash {
	return a.passwordHash
}

func (a *Account) PersonID() uuid.UUID {
	return a.personID
}

func (a *Account) Role() role.Role {
	return a.role
}

func (a *Account) ChangeLogin(login login.Login) {
	a.login = login
}

func (a *Account) ChangePasswordHash(hash passhash.Hash) {
	a.passwordHash = hash
}

func (a *Account) ChangeRole(role role.Role) {
	a.role = role
}
