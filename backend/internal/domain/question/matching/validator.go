package matching

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
)

func validatePairs(pairs []pair.Pair) error {
	if err := validatePairsCount(pairs); err != nil {
		return err
	}
	if err := validatePairsContainsEmpty(pairs); err != nil {
		return err
	}
	if err := validatePairsDuplicates(pairs); err != nil {
		return err
	}
	return nil
}

func validatePairsCount(pairs []pair.Pair) error {
	if len(pairs) < MinPairs {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MinPairs)
	}

	if len(pairs) > MaxPairs {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MaxPairs)
	}

	return nil
}

func validatePairsDuplicates(pairs []pair.Pair) error {
	seenPrompts := make(map[string]struct{}, len(pairs))
	seenMatches := make(map[string]struct{}, len(pairs))

	for i := range pairs {
		promptID := pairs[i].PromptID().String()
		if _, ok := seenPrompts[promptID]; ok {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
		seenPrompts[promptID] = struct{}{}

		matchID := pairs[i].MatchID().String()
		if _, ok := seenMatches[matchID]; ok {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
		seenMatches[matchID] = struct{}{}
	}

	return nil
}

func validatePairsContainsEmpty(pairs []pair.Pair) error {
	for i := range pairs {
		if pairs[i].IsIncomplete() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	return nil
}
