package version

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateTitle(title string) error {
	if err := validateNotEmptyStr(title); err != nil {
		return err
	}

	if err := validateLimitStr(title); err != nil {
		return err
	}

	return nil
}

func validateNotEmptyStr(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w, детали: заголовок версии должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateLimitStr(s string) error {
	if utf8.RuneCountInString(s) > titleCharsLimit {
		return fmt.Errorf("%w, детали: заголовок должен содержать не более %d символов", ErrInvalid, titleCharsLimit)
	}

	return nil
}
