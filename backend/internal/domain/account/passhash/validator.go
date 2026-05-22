package passhash

import (
	"fmt"
)

func validateHash(hash string) error {
	if err := validateRequiredHash(hash); err != nil {
		return err
	}

	if err := validateHashSize(hash); err != nil {
		return err
	}

	return nil
}

func validateRequiredHash(hash string) error {
	if hash == "" {
		return fmt.Errorf("%w: хеш пароля не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateHashSize(hash string) error {
	if len(hash) > MaxHashBytes {
		return fmt.Errorf("%w: хеш не должен весить более %d байт", ErrInvalid, MaxHashBytes)
	}

	return nil
}
