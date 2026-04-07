package account

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Account struct {
	id           uuid.UUID
	login        string
	passwordHash string
	roleID       uuid.UUID
	personID     uuid.UUID
	isBlocked    bool
	createdAt    time.Time
	updatedAt    time.Time
}

func New(login, passwordHash string, roleID, personID uuid.UUID) (*Account, error) {
	if strings.TrimSpace(login) == "" {
		return nil, ErrEmptyLogin
	}

	if strings.TrimSpace(passwordHash) == "" {
		return nil, ErrEmptyPasswordHash
	}

	id, err := uuid.NewV7()

	if err != nil {
		return nil, fmt.Errorf("generate id error: %w", err)
	}

	now := time.Now()

	return &Account{
		id:           id,
		login:        login,
		passwordHash: passwordHash,
		roleID:       roleID,
		personID:     personID,
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
	if strings.TrimSpace(login) == "" {
		return ErrEmptyLogin
	}

	a.login = login
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) ChangePassword(hash string) error {
	if strings.TrimSpace(hash) == "" {
		return ErrEmptyPasswordHash
	}

	a.passwordHash = hash
	a.updatedAt = time.Now()
	return nil
}

func (a *Account) ComparePassword(plain string, comparer PasswordComparer) bool {
	return comparer.Compare(a.passwordHash, plain)
}

func (a *Account) RoleID() uuid.UUID {
	return a.roleID
}

func (a *Account) ChangeRoleID(id uuid.UUID) {
	a.roleID = id
	a.updatedAt = time.Now()
}

func (a *Account) PersonID() uuid.UUID {
	return a.personID
}

func (a *Account) ChangePersonID(id uuid.UUID) {
	a.personID = id
	a.updatedAt = time.Now()
}

func (a *Account) IsBlocked() bool {
	return a.isBlocked
}

func (a *Account) Block() {
	a.isBlocked = true
	a.updatedAt = time.Now()
}

func (a *Account) Unblock() {
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
