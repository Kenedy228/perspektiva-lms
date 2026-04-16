package account

import (
	"time"

	"gitflic.ru/lms/internal/domain/role"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Account struct {
	id           uuid.UUID
	login        string
	passwordHash string
	role         role.Role
	personID     uuid.UUID
	isBlocked    bool
	createdAt    time.Time
	updatedAt    time.Time
}

func New(params Params) (*Account, error) {
	if err := validateLogin(params.Login); err != nil {
		return nil, err
	}

	if err := validatePasswordHash(params.PasswordHash); err != nil {
		return nil, err
	}

	if err := validatePersonID(params.PersonID); err != nil {
		return nil, err
	}

	id, err := utils.GenerateID()

	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Account{
		id:           id,
		login:        params.Login,
		passwordHash: params.PasswordHash,
		role:         params.Role,
		personID:     params.PersonID,
		isBlocked:    false,
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

func (a *Account) ChangeLogin(login string) error {
	if err := validateLogin(login); err != nil {
		return err
	}

	a.login = login
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) ChangePassword(hash string) error {
	if err := validatePasswordHash(hash); err != nil {
		return err
	}

	a.passwordHash = hash
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) ComparePassword(plain string, comparer PasswordComparer) bool {
	return comparer.Compare(a.passwordHash, plain)
}

func (a *Account) Role() role.Role {
	return a.role
}

func (a *Account) ChangeRole(role role.Role) {
	a.role = role
	a.updatedAt = time.Now()
}

func (a *Account) PersonID() uuid.UUID {
	return a.personID
}

func (a *Account) ChangePersonID(id uuid.UUID) error {
	if err := validatePersonID(id); err != nil {
		return err
	}
	a.personID = id
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) IsBlocked() bool {
	return a.isBlocked
}

func (a *Account) Block() {
	if a.IsBlocked() {
		return
	}
	a.isBlocked = true
	a.updatedAt = time.Now()
}

func (a *Account) Unblock() {
	if !a.IsBlocked() {
		return
	}
	a.isBlocked = false
	a.updatedAt = time.Now()
}

func (a *Account) CreatedAt() time.Time {
	return a.createdAt
}

func (a *Account) UpdatedAt() time.Time {
	return a.updatedAt
}

func (a *Account) Equal(other *Account) bool {
	if other == nil {
		return false
	}

	return a.id == other.id
}
