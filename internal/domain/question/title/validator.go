package title

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/internal/domain/question/content"
)

func validateContent(c content.Content) error {
	if err := validateAllowedType(c); err != nil {
		return err
	}

	if err := validateValueLimit(c); err != nil {
		return err
	}

	return nil
}

func validateAllowedType(c content.Content) error {
	if !c.IsText() {
		return fmt.Errorf("%w, детали: заголовок должен быть текстового типа", ErrInvalid)
	}

	return nil
}

func validateValueLimit(c content.Content) error {
	if utf8.RuneCountInString(c.Value()) > titleCharsLimit {
		return fmt.Errorf("%w, детали: заголовок не должен превышать %d символов", ErrInvalid, titleCharsLimit)
	}

	return nil
}
