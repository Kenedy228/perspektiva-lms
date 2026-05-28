package quiz

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
	"gitflic.ru/lms/backend/internal/domain/shared/duplicates"
	"github.com/google/uuid"
)

func validateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор теста обязателен", ErrInvalid)
	}

	return nil
}

func validateTitle(t title.Title) error {
	if t.IsZero() {
		return fmt.Errorf("%w: заголовок теста обязателен", ErrInvalid)
	}

	return nil
}

func validateSources(sources []source.Source) error {
	if len(sources) == 0 {
		return fmt.Errorf("%w: список источников не может быть пустым", ErrInvalid)
	}

	if len(sources) > maxSourcesCount {
		return fmt.Errorf("%w: превышено максимальное количество источников (%d)", ErrInvalid, maxSourcesCount)
	}

	bankIDs := getBankIDs(sources)
	if has := duplicates.HasUUID(bankIDs); has {
		return fmt.Errorf("%w: в тесте обнаружены дублирующиеся банки вопросов", ErrInvalid)
	}

	for i := range sources {
		if sources[i].IsZero() {
			return fmt.Errorf("%w: источник вопросов не должен быть пустым", ErrInvalid)
		}
	}

	return nil
}

func validateSourceToAdd(sources []source.Source, s source.Source) error {
	if s.IsZero() {
		return fmt.Errorf("%w: добавляемый источник не должен быть пустым", ErrInvalid)
	}

	if len(sources) >= maxSourcesCount {
		return fmt.Errorf("%w: превышено максимальное количество источников", ErrInvalid)
	}

	for _, existing := range sources {
		if existing.BankID() == s.BankID() {
			return fmt.Errorf("%w: банк вопросов уже добавлен в тест", ErrInvalid)
		}
	}

	return nil
}

func validateSourceExists(sources []source.Source, bankID uuid.UUID) error {
	for i := range sources {
		if sources[i].BankID() == bankID {
			return nil
		}
	}

	return fmt.Errorf("%w: источник с указанным банком не найден в тесте", ErrInvalid)
}

func validateSourcesToRemove(sources []source.Source) error {
	if len(sources) == 1 {
		return fmt.Errorf("%w: нельзя удалить последний источник вопросов", ErrInvalid)
	}

	return nil
}
