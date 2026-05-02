package attachment

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question/content"
)

func validateContent(c content.Content) error {
	if err := validateAllowedType(c); err != nil {
		return err
	}

	return nil
}

func validateAllowedType(c content.Content) error {
	if c.IsText() {
		return fmt.Errorf("%w, детали: текст не может использоваться в качестве вложения", ErrInvalid)
	}

	return nil
}
