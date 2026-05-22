package account

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/account/passhash"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор аккаунта не существует", ErrInvalid)
	}

	return nil
}

func validatePersonID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор обладателя аккаунта не существует", ErrInvalid)
	}

	return nil
}

func validateLogin(l login.Login) error {
	if l.IsZero() {
		return fmt.Errorf("%w: логин не может быть пустым", ErrInvalid)
	}

	return nil
}

func validatePasswordHash(h passhash.Hash) error {
	if h.IsZero() {
		return fmt.Errorf("%w: хеш пароля не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateRole(r role.Role) error {
	if r.IsZero() || !r.Kind().IsValid() {
		return fmt.Errorf("%w: указана несуществующая или некорректная роль аккаунта", ErrInvalid)
	}

	return nil
}

func validateStatus(s Status) error {
	if !s.IsValid() {
		return fmt.Errorf("%w: указан несуществующий или некорректный статус аккаунта", ErrInvalid)
	}

	return nil
}
