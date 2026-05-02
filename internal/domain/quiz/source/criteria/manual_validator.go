package criteria

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/internal/domain/shared/duplicates"
	"github.com/google/uuid"
)

func validateQuestionIDs(questionIDs []uuid.UUID) error {
	if err := validateQuestionIDsBoundaries(questionIDs); err != nil {
		return err
	}

	if err := validateQuestionIDsNil(questionIDs); err != nil {
		return err
	}

	if err := validateQuestionIDsDuplicates(questionIDs); err != nil {
		return err
	}

	return nil
}

func validateQuestionIDsBoundaries(questionIDs []uuid.UUID) error {
	if len(questionIDs) == 0 {
		return fmt.Errorf("%w, детали: выборка должна содержать хотя бы один вопрос", ErrInvalid)
	}

	if len(questionIDs) > maxQuestionsCount {
		return fmt.Errorf("%w, детали: максимальный размер выборки - %d вопросов", ErrInvalid, maxQuestionsCount)
	}

	return nil
}

func validateQuestionIDsNil(questionIDs []uuid.UUID) error {
	if slices.Contains(questionIDs, uuid.Nil) {
		return fmt.Errorf("%w, детали: один из вопросов в выборке не существует или был удален", ErrInvalid)
	}

	return nil
}

func validateQuestionIDsDuplicates(questionIDs []uuid.UUID) error {
	if has := duplicates.HasUUID(questionIDs); has {
		return fmt.Errorf("%w, детали: выборка содержит дубликаты", ErrInvalid)
	}

	return nil
}
