package account

import (
	"time"

	"gitflic.ru/lms/internal/domain/role"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Account struct {
	id           uuid.UUID
	login        string
	passwordHash string
	role         role.Role
	personID     uuid.UUID
	createdAt    time.Time
	updatedAt    time.Time
}

func New(params Params) (*Account, error) {
	login := normalize(params.Login)

	if err := validateLogin(login); err != nil {
		return nil, err
	}

	if err := validatePasswordHash(params.PasswordHash); err != nil {
		return nil, err
	}

	if err := validatePersonID(params.PersonID); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Account{
		id:           id,
		login:        login,
		passwordHash: params.PasswordHash,
		role:         params.Role,
		personID:     params.PersonID,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

func (a *Account) ID() uuid.UUID {
	return a.id
}

func (a *Account) Login() string {
	return a.login
}

func (a *Account) PasswordHash() string {
	return a.passwordHash
}

func (a *Account) PersonID() uuid.UUID {
	return a.personID
}

func (a *Account) Role() role.Role {
	return a.role
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) ChangeLogin(login string) error {
	nLogin := normalize(login)

	if err := validateLogin(nLogin); err != nil {
		return err
	}

	a.login = nLogin
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) ChangePasswordHash(hash string) error {
	if err := validatePasswordHash(hash); err != nil {
		return err
	}

	a.passwordHash = hash
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) ChangeRole(role role.Role) {
	a.role = role
	a.updatedAt = time.Now()
}

func (a *Account) ChangePersonID(personID uuid.UUID) error {
	if err := validatePersonID(personID); err != nil {
		return err
	}

	a.personID = personID
	a.updatedAt = time.Now()
	return nil
}
