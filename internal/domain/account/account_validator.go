package account

import (
	"errors"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrEmptyLogin        = errors.New("empty login")
	ErrEmptyPasswordHash = errors.New("empty password hash")
	ErrNilPersonID       = errors.New("nil person id")
)

func validateLogin(login string) error {
	if strings.TrimSpace(login) == "" {
		return ErrEmptyLogin
	}

	return nil
}

func validatePasswordHash(passwordHash string) error {
	if strings.TrimSpace(passwordHash) == "" {
		return ErrEmptyPasswordHash
	}

	return nil
}

func validatePersonID(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrNilPersonID
	}

	return nil
}
