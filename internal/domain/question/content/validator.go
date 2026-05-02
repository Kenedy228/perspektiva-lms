package content

import (
	"fmt"
	"strings"
)

func validateValue(value string) error {
	if err := validateRequiredValue(value); err != nil {
		return err
	}

	return nil
}

func validateRequiredValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%w, детали: содержимое контента не может быть пустым", ErrInvalid)
	}

	return nil
}
