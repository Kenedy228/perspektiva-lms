package attempt

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

func validateEnrollmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: ID зачисления не может быть пустым", ErrInvalid)
	}
	return nil
}

func validateQuizID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: ID теста не может быть пустым", ErrInvalid)
	}
	return nil
}

func validateQuestions(questions []question.Question) error {
	if len(questions) == 0 {
		return fmt.Errorf("%w: список вопросов пуст", ErrInvalid)
	}

	for i, q := range questions {
		if q == nil {
			return fmt.Errorf("%w: вопрос под индексом %d равен nil", ErrInvalid, i)
		}
	}
	return nil
}
