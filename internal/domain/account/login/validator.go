package login

import (
	"fmt"
	"strings"
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
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w, детали: логин должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateValueCharsLimit(value string) error {
	if utf8.RuneCountInString(value) < minValueCharsCount {
		return fmt.Errorf("%w, детали: логин должен содержать минимум %d непробельных символов", ErrInvalid, minValueCharsCount)
	}

	if utf8.RuneCountInString(value) > maxValueCharsCount {
		return fmt.Errorf("%w, детали: логин должен содержать максимум %d непробельных символов", ErrInvalid, maxValueCharsCount)
	}

	return nil
}

func validateAllowedCharacters(value string) error {
	if !validValuePattern.MatchString(value) {
		return fmt.Errorf("%w, детали: логин может содержать только латинские буквы, цифры, точки, дефисы и знак подчеркивания", ErrInvalid)
	}

	return nil
}
