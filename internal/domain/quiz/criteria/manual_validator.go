package criteria

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/internal/domain/shared/duplicate"
	"github.com/google/uuid"
)

func validateQuestionIDs(questionIDs []uuid.UUID) error {
	if len(questionIDs) == 0 {
		return fmt.Errorf("%w, детали: выборка должна содержать хотя бы один вопрос", ErrInvalidQuestions)
	}

	if len(questionIDs) > maxQuestions {
		return fmt.Errorf("%w, детали: максимальный размер выборки - %d вопросов", ErrInvalidQuestions, maxQuestions)
	}

	if slices.Contains(questionIDs, uuid.Nil) {
		return fmt.Errorf("%w, детали: один из вопросов в выборке не существует или был удален", ErrInvalidQuestions)
	}

	if has := duplicate.FindUUID(questionIDs); has {
		return fmt.Errorf("%w, детали: выборка содержит дубликаты", ErrInvalidQuestions)
	}

	return nil
}
