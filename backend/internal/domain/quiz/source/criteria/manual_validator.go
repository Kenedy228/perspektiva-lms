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
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if len(questionIDs) > maxQuestionsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, maxQuestionsCount)
	}

	return nil
}

func validateQuestionIDsNil(questionIDs []uuid.UUID) error {
	if slices.Contains(questionIDs, uuid.Nil) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateQuestionIDsDuplicates(questionIDs []uuid.UUID) error {
	if has := duplicates.HasUUID(questionIDs); has {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
