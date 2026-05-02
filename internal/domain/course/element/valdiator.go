package element

import (
	"fmt"
)

func validateContent(content Content) error {
	if err := validateRequiredContent(content); err != nil {
		return err
	}

	return nil
}

func validateRequiredContent(c Content) error {
	if c == nil {
		return fmt.Errorf("%w, детали: элемент курса не может существовать без наполнения", ErrInvalid)
	}

	return nil
}
