package answer

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/selectable"
)

func validateOptionIDs(optionIDs []OptionID) error {
	if err := validateOptionIDsCount(optionIDs); err != nil {
		return err
	}

	if err := validateOptionIDsContainsEmpty(optionIDs); err != nil {
		return err
	}

	if err := validateOptionIDsDuplicates(optionIDs); err != nil {
		return err
	}

	return nil
}

func validateOptionIDsCount(optionIDs []OptionID) error {
	if len(optionIDs) > selectable.MaxOptionsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, selectable.MaxOptionsCount)
	}

	return nil
}

func validateOptionIDsContainsEmpty(optionIDs []OptionID) error {
	for i := range optionIDs {
		if optionIDs[i].IsZero() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	return nil
}

func validateOptionIDsDuplicates(optionIDs []OptionID) error {
	for i := range optionIDs {
		for j := i + 1; j < len(optionIDs); j++ {
			if optionIDs[i].ID() == optionIDs[j].ID() {
				return fmt.Errorf("%w: invalid value", ErrInvalid)
			}
		}
	}

	return nil
}
