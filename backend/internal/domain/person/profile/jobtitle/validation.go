package jobtitle

import (
	"fmt"
	"unicode/utf8"
)

func validateValue(value string) error {
	if value == "" {
		return nil
	}

	if err := validateValueCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateValueCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)

	if rc > ValueCharsLimit {
		return fmt.Errorf("%w: сведение о занимаемой должности не может превышать %d символов (текущее количество символов - %d штук)", ErrInvalid, ValueCharsLimit, rc)
	}

	return nil
}
