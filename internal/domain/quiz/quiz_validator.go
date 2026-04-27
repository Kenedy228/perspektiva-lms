package quiz

import (
	"fmt"
	"slices"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/duplicate"
)

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w, детали: заголовок должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateSources(sources []Source) error {
	if len(sources) == 0 {
		return fmt.Errorf("%w, детали: квиз должен иметь хотя бы один источник вопросов", ErrInvalid)
	}

	if len(sources) > maxSources {
		return fmt.Errorf("%w, детали: квиз должен содержать не более %d источников", ErrInvalid, maxSources)
	}

	bankIDs := getBankIDs(sources)
	if has := duplicate.FindUUID(bankIDs); has {
		return fmt.Errorf("%w, детали: квиз не может содержать источники с одинаковым банком вопросов", ErrInvalid)
	}

	return nil
}

func validateMaxAttempts(maxAttempts int) error {
	if maxAttempts < 0 {
		return fmt.Errorf("%w, детали: максимальное число попыток должно быть неотрицательным", ErrInvalid)
	}

	if maxAttempts > maxAttemptsCount {
		return fmt.Errorf("%w, детали: максимальное число попыток должно быть менее %d", ErrInvalid, maxAttemptsCount)
	}

	return nil
}

func validateSourceToAdd(sources []Source, s Source) error {
	if len(sources)+1 > maxSources {
		return ErrSourceSizeExceeded
	}

	cSources := slices.Clone(sources)
	cSources = append(cSources, s)
	bankIDs := getBankIDs(cSources)

	if has := duplicate.FindUUID(bankIDs); has {
		return fmt.Errorf("%w, детали: источник с таким банком уже указан в квизе", ErrDuplicateBankID)
	}

	return nil
}
