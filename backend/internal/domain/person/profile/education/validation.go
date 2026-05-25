package education

import (
	"fmt"
	"unicode/utf8"
)

// validateValue валидирует сведения об образовании.
func validateValue(value string) error {
	if err := validateRequiredValue(value); err != nil {
		return err
	}

	if err := validateValueCharsLimit(value); err != nil {
		return err
	}

	return nil
}

// validateRequiredValue проверяет, что переданные сведения об образовании не являются пустыми.
func validateRequiredValue(value string) error {
	if value == "" {
		return fmt.Errorf("%w: сведения об образовании не могут быть пустыми", ErrInvalid)
	}

	return nil
}

// validateValueCharsLimit проверяет, что количество символов в сведениях об образовании не превышает лимита.
func validateValueCharsLimit(value string) error {
	if rc := utf8.RuneCountInString(value); rc > ValueCharsLimit {
		return fmt.Errorf("%w: сведения об образовании не могут превышать %d символов (получено символов %d)", ErrInvalid, ValueCharsLimit, rc)
	}

	return nil
}
