package answer

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

func validateQuestionID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: идентификатор вопроса не может быть пустым", ErrInvalid)
	}
	return nil
}

func validateAnswer(a question.Answer) error {
	if a == nil {
		return fmt.Errorf("%w, детали: ответ не может быть пустым", ErrInvalid)
	}
	return nil
}
