package passhash

import (
	"fmt"
	"strings"
)

func validateHash(hash string) error {
	if err := validateRequiredHash(hash); err != nil {
		return err
	}

	if err := validateHashFormat(hash); err != nil {
		return err
	}

	return nil
}

func validateRequiredHash(hash string) error {
	if strings.TrimSpace(hash) == "" {
		return fmt.Errorf("%w, детали: хеш пароля не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateHashFormat(hash string) error {
	if !bcryptFormat.MatchString(hash) {
		return fmt.Errorf("%w, детали: неверный формат криптографического хеша", ErrInvalid)
	}

	return nil
}
