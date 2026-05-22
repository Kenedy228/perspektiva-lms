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

func validateSources(sources []source.Source) error {
	if len(sources) == 0 {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if len(sources) > maxSourcesCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, maxSourcesCount)
	}

	bankIDs := getBankIDs(sources)
	if has := duplicates.HasUUID(bankIDs); has {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	for i := range sources {
		if sources[i].IsZero() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	return nil
}

func validateSourceToAdd(sources []source.Source, s source.Source) error {
	if s.IsZero() {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if len(sources) >= maxSourcesCount {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	for _, existing := range sources {
		if existing.BankID() == s.BankID() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
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

	return fmt.Errorf("%w: invalid value", ErrInvalid)
}

func validateSourcesToRemove(sources []source.Source) error {
	if len(sources) == 1 {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	return nil
}
