package account

import (
	"fmt"

	"github.com/google/uuid"
)

func validateLogin(login string) error {
	if err := validateRequired(loginField, login); err != nil {
		return err
	}

	if err := validateLength(loginField, login, loginCharsLimit); err != nil {
		return err
	}

	return nil
}

func validatePasswordHash(hash string) error {
	if err := validateRequired(passwordField, hash); err != nil {
		return err
	}

	return nil
}

func validatePersonID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: детали человека обязательны для заполнения", ErrInvalid)
	}

	return nil
}
