package option

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
)

func validateText(t text.Text) error {
	if err := validateTextCharsLimit(t); err != nil {
		return err
	}

	return nil
}

func validateTextCharsLimit(t text.Text) error {
	if utf8.RuneCountInString(t.Value()) > TextCharsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, TextCharsLimit)
	}

	return nil
}
