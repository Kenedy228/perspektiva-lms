package education

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateRequiredValue(value); err != nil {
		return err
	}

	if err := validateCharsLimitEducation(value); err != nil {
		return err
	}

	return nil
}

func validateRequiredValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w, детали: образование должно содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateCharsLimitEducation(value string) error {
	if utf8.RuneCountInString(value) > ValueCharsLimit {
		return fmt.Errorf("%w, детали: образование должно содержать не более %d символов", ErrInvalid, ValueCharsLimit)
	}

	return nil
}
