package login

import (
	"fmt"
	"unicode/utf8"
)

func validateValue(value string) error {
	if err := validateRequiredValue(value); err != nil {
		return err
	}

	if err := validateValueCharsLimit(value); err != nil {
		return err
	}

	if err := validateAllowedCharacters(value); err != nil {
		return err
	}

	return nil
}

func validateRequiredValue(value string) error {
	if value == "" {
		return fmt.Errorf("%w: логин не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateValueCharsLimit(value string) error {
	if utf8.RuneCountInString(value) < MinValueCharsCount {
		return fmt.Errorf("%w: логин должен состоять не менее чем из %d символов", ErrInvalid, MinValueCharsCount)
	}

	if utf8.RuneCountInString(value) > MaxValueCharsCount {
		return fmt.Errorf("%w: логин должен состоять не более чем из %d символов", ErrInvalid, MaxValueCharsCount)
	}

	return nil
}

func validateAllowedCharacters(value string) error {
	if !validValuePattern.MatchString(value) {
		return fmt.Errorf("%w: логин содержит запрещенные символы", ErrInvalid)
	}

	return nil
}
