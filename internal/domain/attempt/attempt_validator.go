package attempt

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

func validateEnrollmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: зачисление не существует (nil)", ErrInvalidAttempt)
	}

	return nil
}

func validateQuizID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: тест не существует (nil)", ErrInvalidAttempt)
	}

	return nil
}

func validateQuestions(questions []question.Question) error {
	if len(questions) == 0 {
		return fmt.Errorf("%w", ErrInvalidAttempt)
	}

	for i := range questions {
		if questions[i] == nil {
			return fmt.Errorf("%w", ErrInvalidAttempt)
		}
	}

	return nil
}
