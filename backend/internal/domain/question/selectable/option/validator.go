package option

import (
	"fmt"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateRequired(value); err != nil {
		return err
	}

	if err := validateCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateRequired(value string) error {
	if value == "" {
		return fmt.Errorf("%w: значение не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)

	if rc > ValueCharsLimit {
		return fmt.Errorf("%w: значение не может прешывать %d символов (текущее количество символов - %d)", ErrInvalid, ValueCharsLimit, rc)
	}

	return nil
}
