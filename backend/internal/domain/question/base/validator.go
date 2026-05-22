package base

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"github.com/google/uuid"
)

func validateTitle(t title.Title) error {
	if err := validateTitleRequired(t); err != nil {
		return err
	}

	return nil
}

func validateTitleRequired(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: заголовок вопроса не может быть пустым", ErrInvalid)
	}

	return nil
}

func validateID(id uuid.UUID) error {
	if err := validateIDRequired(id); err != nil {
		return err
	}

	return nil
}

func validateIDRequired(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор не может быть пустым", ErrInvalid)
	}

	return nil
}
