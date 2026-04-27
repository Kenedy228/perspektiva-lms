package name

import "fmt"

func validateRequiredPart(field, value string) error {
	if value == "" {
		return fmt.Errorf("%w, детали: %s не может быть пустым", ErrInvalid, field)
	}

	return validatePartStructure(field, value)
}

func validateOptionalPart(field, value string) error {
	if value == "" {
		return nil
	}

	return validatePartStructure(field, value)
}
