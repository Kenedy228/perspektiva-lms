package variant

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
)

func validateText(t text.Text) error {
	if utf8.RuneCountInString(t.Value()) > TextCharsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, TextCharsLimit)
	}

	return nil
}
