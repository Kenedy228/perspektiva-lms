package text

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateNotEmptyValue(value); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	if err := validateCharsLimit(value); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalid, err)
	}

	return nil
}

func validateNotEmptyValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("значение не может быть пустым")
	}

	return nil
}

func validateCharsLimit(value string) error {
	if utf8.RuneCountInString(value) > ValueCharsLimit {
		return fmt.Errorf("превышен лимит символов (%d)", ValueCharsLimit)
	}

	return nil
}
