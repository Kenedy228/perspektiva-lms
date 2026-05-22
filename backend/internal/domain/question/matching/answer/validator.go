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
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, matching.MaxPairs)
	}

	return nil
}

func validatePairIDsNotEmpty(pairs []Pair) error {
	for i := range pairs {
		if pairs[i].PromptID.ID() == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}

		if pairs[i].MatchID.ID() == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	return nil
}

func validatePairsDuplicates(pairs []Pair) error {
	for i := range pairs {
		for j := i + 1; j < len(pairs); j++ {
			if pairs[i].PromptID == pairs[j].PromptID {
				return fmt.Errorf("%w: invalid value", ErrInvalid)
			}

			if pairs[i].MatchID == pairs[j].MatchID {
				return fmt.Errorf("%w: invalid value", ErrInvalid)
			}
		}
	}

	return nil
}
