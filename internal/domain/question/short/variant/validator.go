package variant

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/internal/domain/question/content"
)

func validateContent(c content.Content) error {
	if err := validateContentAllowedFormat(c); err != nil {
		return err
	}

	if err := validateContentAsTextCharsLimit(c); err != nil {
		return err
	}

	return nil
}

func validateContentAllowedFormat(c content.Content) error {
	if !c.IsText() {
		return fmt.Errorf("%w, детали: вариант может быть указан только в текстовом формате", ErrInvalid)
	}

	return nil
}

func validateContentAsTextCharsLimit(c content.Content) error {
	if utf8.RuneCountInString(c.Value()) > contentAsTextCharsLimit {
		return fmt.Errorf("%w, детали: вариант в текстовом формате должен содержать не более %d символов", ErrInvalid, contentAsTextCharsLimit)
	}

	return nil
}
