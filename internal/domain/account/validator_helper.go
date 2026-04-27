package account

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateLength(field, value string, maxLength int) error {
	if utf8.RuneCountInString(value) > maxLength {
		return fmt.Errorf("%w, детали: %s должен быть менее %d символов", ErrInvalid, field, maxLength)
	}

	return nil
}

func validateRequired(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w, детали: %s должен содержать хотя бы один непробельный символ", ErrInvalid, field)
	}

	return nil
}
