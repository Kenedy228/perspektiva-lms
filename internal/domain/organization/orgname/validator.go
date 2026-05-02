package orgname

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateNotEmpty(value); err != nil {
		return err
	}

	if err := validateCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateNotEmpty(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w, детали: название организации должно содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateCharsLimit(value string) error {
	if utf8.RuneCountInString(value) > valueCharsLimit {
		return fmt.Errorf("%w, детали: название организации не должно содержать более %d символов", ErrInvalid, valueCharsLimit)
	}

	return nil
}
