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
		return fmt.Errorf("%w: идентификатор банка вопросов обязателен", ErrInvalid)
	}
	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: заголовок банка вопросов обязателен", ErrInvalid)
	}
	return nil
}

func validateQuestions(questions []uuid.UUID) error {
	if len(questions) > MaxQuestionsCount {
		return fmt.Errorf("%w: банк не может содержать более %d вопросов", ErrInvalid, MaxQuestionsCount)
	}

	for i := range questions {
		if questions[i] == uuid.Nil {
			return fmt.Errorf("%w: идентификатор вопроса не может быть пустым", ErrInvalid)
		}
	}

	if has := duplicates.HasUUID(questions); has {
		return fmt.Errorf("%w: банк не должен содержать повторяющихся вопросов", ErrInvalid)
	}

	return nil
}

func validateQuestionsForAdding(questions []uuid.UUID, add []uuid.UUID) error {
	totalLen := len(questions) + len(add)

	if totalLen > MaxQuestionsCount {
		return fmt.Errorf("%w: банк не может содержать более %d вопросов", ErrInvalid, MaxQuestionsCount)
	}

	for i := range add {
		if add[i] == uuid.Nil {
			return fmt.Errorf("%w: идентификатор добавляемого вопроса не может быть пустым", ErrInvalid)
		}
	}

	cQuestions := make([]uuid.UUID, 0, totalLen)
	cQuestions = append(cQuestions, questions...)
	cQuestions = append(cQuestions, add...)

	if has := duplicates.HasUUID(cQuestions); has {
		return fmt.Errorf("%w: вопрос уже добавлен в банк", ErrInvalid)
	}

	return nil
}

func validateQuestionsForRemoving(questions []uuid.UUID, remove []uuid.UUID) error {
	if len(remove) == 0 {
		return nil
	}

	for i := range remove {
		if remove[i] == uuid.Nil {
			return fmt.Errorf("%w: идентификатор удаляемого вопроса не может быть пустым", ErrInvalid)
		}
	}

	if has := duplicates.HasUUID(remove); has {
		return fmt.Errorf("%w: список удаляемых вопросов не должен содержать дубликатов", ErrInvalid)
	}

	for i := range remove {
		if !slices.Contains(questions, remove[i]) {
			return fmt.Errorf("%w: вопрос %s не найден в банке", ErrInvalid, remove[i])
		}
	}

	return nil
}
