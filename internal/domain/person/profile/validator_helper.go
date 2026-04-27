package profile

import (
	"fmt"
	"unicode/utf8"
)

func validateMaxLength(field, value string, length int) error {
	if utf8.RuneCountInString(value) > length {
		return fmt.Errorf("%w, детали: поле %q не должно превышать %d символов", ErrInvalid, field, length)
	}

	return nil
}
