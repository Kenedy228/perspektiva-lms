package title

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateValueNotEmpty(value); err != nil {
		return err
	}

	if err := validateValueLimit(value); err != nil {
		return err
	}

	return nil
}

func validateValueNotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateValueLimit(value string) error {
	if utf8.RuneCountInString(value) > ValueCharsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, ValueCharsLimit)
	}

	return nil
}
