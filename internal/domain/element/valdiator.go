package element

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateTitle(title string) error {
	if err := requireNotEmpty(titleField, title); err != nil {
		return err
	}

	if err := requireCharsLimit(titleField, title, titleCharsLimit); err != nil {
		return err
	}

	return nil
}

func validateContent(content Content) error {
	if err := requireNotNil(contentField, content); err != nil {
		return err
	}

	return nil
}

func requireNotEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w, детали: %q должно содержать хотя бы один непробельный символ", ErrInvalid, field)
	}

	return nil
}

func requireCharsLimit(field, value string, maxLen int) error {
	if utf8.RuneCountInString(value) > maxLen {
		return fmt.Errorf("%w, детали: %q не должно превышать %d символов", ErrInvalid, field, maxLen)
	}

	return nil
}

func requireNotNil(field string, value any) error {
	if value == nil {
		return fmt.Errorf("%w, детали: %q должен быть предоставлен", ErrInvalid, field)
	}

	return nil
}
