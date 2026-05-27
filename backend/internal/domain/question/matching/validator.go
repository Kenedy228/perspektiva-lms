package matching

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
)

func validateBase(b *base.Base) error {
	if err := validateRequiredBase(b); err != nil {
		return err
	}

	return nil
}

func validateRequiredBase(b *base.Base) error {
	if b == nil {
		return fmt.Errorf("%w: база вопросов обязательна", ErrInvalid)
	}

	return nil
}

func validatePairs(pairs []pair.Pair) error {
	if err := validatePairsCount(pairs); err != nil {
		return err
	}
	if err := validatePairsContainsEmpty(pairs); err != nil {
		return err
	}
	return nil
}

func validatePairsCount(pairs []pair.Pair) error {
	if len(pairs) < MinPairs {
		return fmt.Errorf(
			"%w: количество пар должно быть не меньше %d (текущее количество: %d)",
			ErrInvalid,
			MinPairs,
			len(pairs),
		)
	}

	if len(pairs) > MaxPairs {
		return fmt.Errorf(
			"%w: количество пар должно быть не больше %d (текущее количество: %d)",
			ErrInvalid,
			MaxPairs,
			len(pairs),
		)
	}

	return nil
}

func validatePairsContainsEmpty(pairs []pair.Pair) error {
	for i := range pairs {
		if pairs[i].IsZero() {
			return fmt.Errorf("%w: пара под индексом %d не заполнена", ErrInvalid, i)
		}
	}

	return nil
}
