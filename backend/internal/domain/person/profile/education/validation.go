package education

import (
	"fmt"
	"unicode/utf8"
)

func validateValue(value string) error {
	if value == "" {
		return nil
	}

	if err := validateCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)
	if rc > ValueCharsLimit {
		return fmt.Errorf("%w: сведение об образовании не может превышать %d символов (текущее количество символов %d штук)", ErrInvalid, ValueCharsLimit, rc)
	}

	return nil
}
