package account

import (
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

// Account представляет собой учетные данные пользователя системы.
type Account struct {
	id           uuid.UUID
	login        login.Login
	passwordHash passhash.Hash
	role         role.Role
	personID     uuid.UUID
	status       Status
}

// New создает новый аккаунт.
func New(l login.Login, h passhash.Hash, r role.Role, personID uuid.UUID) (*Account, error) {
	if err := validateLogin(l); err != nil {
		return nil, err
	}

	if err := validatePasswordHash(h); err != nil {
		return nil, err
	}

	if err := validateRole(r); err != nil {
		return nil, err
	}

	if err := validatePersonID(personID); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Account{
		id:           id,
		login:        l,
		passwordHash: h,
		role:         r,
		personID:     personID,
		status:       StatusActive,
	}, nil
}

// Restore восстанавливает существующий аккаунт.
func Restore(id uuid.UUID, l login.Login, h passhash.Hash, r role.Role, personID uuid.UUID, status Status) (*Account, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := validateLogin(l); err != nil {
		return nil, err
	}

	if err := validatePasswordHash(h); err != nil {
		return nil, err
	}

	if err := validateRole(r); err != nil {
		return nil, err
	}

	if err := validatePersonID(personID); err != nil {
		return nil, err
	}

	if err := validateStatus(status); err != nil {
		return nil, err
	}

	return &Account{
		id:           id,
		login:        l,
		passwordHash: h,
		role:         r,
		personID:     personID,
		status:       status,
	}, nil
}

// ID возвращает идентификатор аккаунта.
func (a *Account) ID() uuid.UUID {
	return a.id
}

// Login возвращает логин аккаунта.
func (a *Account) Login() login.Login {
	return a.login
}

// PasswordHash возвращает хеш пароля аккаунта.
func (a *Account) PasswordHash() passhash.Hash {
	return a.passwordHash
}

// PersonID возвращает идентификатор обладателя аккаунта.
func (a *Account) PersonID() uuid.UUID {
	return a.personID
}

// Role возвращает роль аккаунта.
func (a *Account) Role() role.Role {
	return a.role
}

// Status возвращает статус аккаунта.
func (a *Account) Status() Status {
	return a.status
}

// IsActive сигнализирует, можно ли использовать данный аккаунт при аутентификации.
func (a *Account) IsActive() bool {
	return a.status == StatusActive
}

// IsDeleted сигнализирует, удален ли аккаунт.
func (a *Account) IsDeleted() bool {
	return a.status == StatusDeleted
}

// ChangeLogin изменяет логин аккаунта.
func (a *Account) ChangeLogin(login login.Login) error {
	if err := ensureNotDeleted(a); err != nil {
		return err
	}

	if err := validateLogin(login); err != nil {
		return err
	}

	a.login = login
	return nil
}

// ChangePasswordHash изменяет хеш пароля.
func (a *Account) ChangePasswordHash(hash passhash.Hash) error {
	if err := ensureNotDeleted(a); err != nil {
		return err
	}

	if err := validatePasswordHash(hash); err != nil {
		return err
	}

	a.passwordHash = hash
	return nil
}

// ChangeRole изменяет роль аккаунта.
func (a *Account) ChangeRole(role role.Role) error {
	if err := ensureNotDeleted(a); err != nil {
		return err
	}

	if err := validateRole(role); err != nil {
		return err
	}

	a.role = role
	return nil
}

// Delete помечает аккаунт как удаленный.
func (a *Account) Delete() error {
	if err := ensureNotDeleted(a); err != nil {
		return err
	}

	a.status = StatusDeleted
	return nil
}

// Block блокирует аккаунт (помечает как удаленный).
func (a *Account) Block() error {
	return a.Delete()
}

// Activate восстанавливает аккаунт из удаленного состояния.
func (a *Account) Activate() error {
	if !a.IsDeleted() {
		return ErrAlreadyActive
	}

	a.status = StatusActive
	return nil
}
