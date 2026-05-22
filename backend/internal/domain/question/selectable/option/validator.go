package option

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
)

func validateOptionText(t text.Text) error {
	if err := validateOptionTextCharsLimit(t); err != nil {
		return err
	}

	return nil
}

func validateOptionTextCharsLimit(t text.Text) error {
	if utf8.RuneCountInString(t.Value()) > TextCharsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, TextCharsLimit)
	}

	return nil
}
