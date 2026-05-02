package bank

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/shared/duplicates"
	"github.com/google/uuid"
)

func validateQuestionsForAdding(questions []uuid.UUID, add []uuid.UUID) error {
	totalLen := len(questions) + len(add)

	if totalLen > maxQuestionsCount {
		return fmt.Errorf("%w, детали: количество вопросов превысит лимит %d вопросов", ErrInvalid, maxQuestionsCount)
	}

	for i := range add {
		if add[i] == uuid.Nil {
			return fmt.Errorf("%w, детали: добавляемые вопросы содержат несуществующий вопрос", ErrInvalid)
		}
	}

	cQuestions := make([]uuid.UUID, 0, totalLen)
	cQuestions = append(cQuestions, questions...)
	cQuestions = append(cQuestions, add...)

	if has := duplicates.HasUUID(cQuestions); has {
		return fmt.Errorf("%w, детали: один из добавляемых вопросов либо уже существует в банке вопросов, либо дублирует вопрос в добавляемом списке", ErrInvalid)
	}

	return nil
}
