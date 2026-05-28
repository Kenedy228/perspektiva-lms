package criteria

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/shared/duplicates"
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
		return fmt.Errorf("%w: список идентификаторов вопросов не может быть пустым", ErrInvalid)
	}

	if len(questionIDs) > maxQuestionsCount {
		return fmt.Errorf("%w: количество вопросов не должно превышать %d", ErrInvalid, maxQuestionsCount)
	}

	return nil
}

func validateQuestionIDsNil(questionIDs []uuid.UUID) error {
	if slices.Contains(questionIDs, uuid.Nil) {
		return fmt.Errorf("%w: список не должен содержать пустых идентификаторов вопросов", ErrInvalid)
	}

	return nil
}

func validateQuestionIDsDuplicates(questionIDs []uuid.UUID) error {
	if has := duplicates.HasUUID(questionIDs); has {
		return fmt.Errorf("%w: список не должен содержать повторяющихся идентификаторов вопросов", ErrInvalid)
	}

	return nil
}
