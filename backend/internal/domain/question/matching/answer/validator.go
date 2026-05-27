package answer

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"github.com/google/uuid"
)

func validateAnswerPairs(pairs []Pair) error {
	if err := validateAnswerPairsMaxCount(pairs); err != nil {
		return err
	}

	if err := validatePairIDsNotEmpty(pairs); err != nil {
		return err
	}

	if err := validatePairsDuplicates(pairs); err != nil {
		return err
	}

	return nil
}

func validateAnswerPairsMaxCount(pairs []Pair) error {
	if len(pairs) > matching.MaxPairs {
		return fmt.Errorf(
			"%w: количество пар не должно превышать %d штук (текущее количество - %d)",
			ErrInvalid,
			matching.MaxPairs,
			len(pairs),
		)
	}

	return nil
}

func validatePairIDsNotEmpty(pairs []Pair) error {
	for i := range pairs {
		if pairs[i].PromptID == uuid.Nil {
			return fmt.Errorf("%w: в паре с индексом %d отсутствует идентификатор prompt", ErrInvalid, i)
		}

		if pairs[i].MatchID == uuid.Nil {
			return fmt.Errorf("%w: в паре с индексом %d отсутствует идентификатор match", ErrInvalid, i)
		}
	}

	return nil
}

func validatePairsDuplicates(pairs []Pair) error {
	for i := range pairs {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].PromptID == pairs[j].PromptID {
				return fmt.Errorf(
					"%w: prompt %s повторяется в парах с индексами %d и %d",
					ErrInvalid,
					pairs[i].PromptID,
					i,
					j,
				)
			}

			if pairs[i].MatchID == pairs[j].MatchID {
				return fmt.Errorf(
					"%w: match %s повторяется в парах с индексами %d и %d",
					ErrInvalid,
					pairs[i].MatchID,
					i,
					j,
				)
			}
		}
	}

	return nil
}
