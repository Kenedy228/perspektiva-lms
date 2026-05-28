package title

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateNotEmpty(value); err != nil {
		return err
	}

	if err := validateValueCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateNotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: название блока не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateValueCharsLimit(value string) error {
	if utf8.RuneCountInString(value) > valueCharsLimit {
		return fmt.Errorf("%w: название блока не должно превышать %d символов", ErrInvalid, valueCharsLimit)
	}

	return nil
}
