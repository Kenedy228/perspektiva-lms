package matching

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question/matching/pair"
)

func validatePairs(pairs []pair.Pair) error {
	if err := validatePairsCount(pairs); err != nil {
		return err
	}

	return nil
}

func validatePairsCount(pairs []pair.Pair) error {
	if len(pairs) < minPairs {
		return fmt.Errorf("%w, детали: вопрос должен содержать минимум %d пар соответствий", ErrInvalid, minPairs)
	}

	if len(pairs) > maxPairs {
		return fmt.Errorf("%w, детали: вопрос должен содержать максимум %d пар соответствий", ErrInvalid, maxPairs)
	}

	return nil
}
