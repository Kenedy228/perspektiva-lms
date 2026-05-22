package name

import (
	"fmt"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateRequired(value); err != nil {
		return err
	}

	if err := validateValueCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateRequired(value string) error {
	if value == "" {
		return fmt.Errorf("%w: наименование не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateValueCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)

	if rc > ValueCharsLimit {
		return fmt.Errorf("%w: наименование должно содержать не более %d символов, действительно содержит %d символов", ErrInvalid, ValueCharsLimit, rc)
	}

	return nil
}
