package item

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question"
)

func validateQuestion(q question.Question) error {
	if err := validateRequired(q); err != nil {
		return err
	}

	return nil
}

func validateRequired(q question.Question) error {
	if q == nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
