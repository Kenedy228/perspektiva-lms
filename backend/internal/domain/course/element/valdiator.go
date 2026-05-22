package element

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/course/element/title"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateContent(content Content) error {
	if err := validateRequiredContent(content); err != nil {
		return err
	}

	return nil
}

func validateRequiredContent(c Content) error {
	if c == nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
