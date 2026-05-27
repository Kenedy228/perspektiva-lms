package option

import (
	"fmt"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateValueRequired(value); err != nil {
		return err
	}

	if err := validateValueCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateValueRequired(value string) error {
	if value == "" {
		return fmt.Errorf("%w: значение опции не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateValueCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)

	if rc > ValueCharsLimit {
		return fmt.Errorf("%w: значение опции не должно превышать %d символов (текущее количество символов - %d)", ErrInvalid, ValueCharsLimit, rc)
	}

	return nil
}
