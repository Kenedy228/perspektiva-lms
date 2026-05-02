package item

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
)

func validateQuestion(q question.Question) error {
	if err := validateRequired(q); err != nil {
		return err
	}

	return nil
}

func validateRequired(q question.Question) error {
	if q == nil {
		return fmt.Errorf("%w, детали: переданный вопрос не существует", ErrInvalid)
	}

	return nil
}
