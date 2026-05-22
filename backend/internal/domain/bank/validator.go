package bank

import (
	"fmt"
	"slices"

	"gitflic.ru/lms/backend/internal/domain/bank/title"
	"gitflic.ru/lms/backend/internal/domain/shared/duplicates"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	return nil
}

func validateQuestions(questions []uuid.UUID) error {
	if len(questions) > maxQuestionsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, maxQuestionsCount)
	}

	for i := range questions {
		if questions[i] == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	if has := duplicates.HasUUID(questions); has {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateQuestionsForAdding(questions []uuid.UUID, add []uuid.UUID) error {
	totalLen := len(questions) + len(add)

	if totalLen > maxQuestionsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, maxQuestionsCount)
	}

	for i := range add {
		if add[i] == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	cQuestions := make([]uuid.UUID, 0, totalLen)
	cQuestions = append(cQuestions, questions...)
	cQuestions = append(cQuestions, add...)

	if has := duplicates.HasUUID(cQuestions); has {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}

func validateQuestionsForRemoving(questions []uuid.UUID, remove []uuid.UUID) error {
	if len(remove) == 0 {
		return nil
	}

	for i := range remove {
		if remove[i] == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	if has := duplicates.HasUUID(remove); has {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	for i := range remove {
		if !slices.Contains(questions, remove[i]) {
			return fmt.Errorf("%w: invalid value (%s)", ErrInvalid, remove[i])
		}
	}

	return nil
}
