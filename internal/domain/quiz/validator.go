package quiz

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/quiz/source"
	"gitflic.ru/lms/internal/domain/shared/duplicates"
)

func validateSources(sources []source.Source) error {
	if len(sources) == 0 {
		return fmt.Errorf("%w, детали: квиз должен иметь хотя бы один источник вопросов", ErrInvalid)
	}

	if len(sources) > maxSourcesCount {
		return fmt.Errorf("%w, детали: квиз должен содержать не более %d источников", ErrInvalid, maxSourcesCount)
	}

	bankIDs := getBankIDs(sources)
	if has := duplicates.HasUUID(bankIDs); has {
		return fmt.Errorf("%w, детали: квиз не может содержать источники с одинаковым банком вопросов", ErrInvalid)
	}

	return nil
}

func validateSourceToAdd(sources []source.Source, s source.Source) error {
	if len(sources) >= maxSourcesCount {
		return fmt.Errorf("%w, детали: достигнуто максимальное количество источников", ErrInvalid)
	}

	for _, existing := range sources {
		if existing.BankID() == s.BankID() {
			return fmt.Errorf("%w, детали: источник с таким банком уже указан", ErrInvalid)
		}
	}

	return nil
}

func validateSourcesToRemove(sources []source.Source) error {
	if len(sources) == 1 {
		return fmt.Errorf("%w, детали: нельзя удалить последний источник в квизе", ErrInvalid)
	}

	return nil
}
