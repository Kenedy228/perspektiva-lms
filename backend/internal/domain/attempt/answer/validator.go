package answer

import (
	"fmt"
	"time"

	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

func validateQuestionID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateAnswer(a question.Answer) error {
	if a == nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateAnsweredAt(answeredAt time.Time) error {
	if answeredAt.IsZero() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}
