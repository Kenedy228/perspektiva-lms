package pair

import (
	"fmt"
	"unicode/utf8"

	"github.com/google/uuid"
)

func validateValue(value string) error {
	if err := validateRequired(value); err != nil {
		return err
	}

	if err := validateCharsLimit(value); err != nil {
		return err
	}

	return nil
}

func validateRequired(value string) error {
	if value == "" {
		return fmt.Errorf("%w: значение не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateCharsLimit(value string) error {
	rc := utf8.RuneCountInString(value)

	if rc > ValueCharsLimit {
		return fmt.Errorf(
			"%w: значение не может прешывать %d символов (текущее количество символов - %d)",
			ErrInvalid,
			ValueCharsLimit,
			rc,
		)
	}

	return nil
}

func validateIDRequired(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор не может быть пустым", ErrInvalid)
	}

	return nil
}

func validatePrompt(prompt Prompt) error {
	if prompt.IsZero() {
		return fmt.Errorf("%w: prompt не должен быть zero-value", ErrInvalid)
	}

	return nil
}

func validateMatch(match Match) error {
	if match.IsZero() {
		return fmt.Errorf("%w: match не должен быть zero-value", ErrInvalid)
	}

	return nil
}
