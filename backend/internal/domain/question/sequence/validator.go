package sequence

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"github.com/google/uuid"
)

func validateOptions(options []option.Option) error {
	if err := validateOptionsCount(options); err != nil {
		return err
	}

	if err := validateOptionsContainsEmpty(options); err != nil {
		return err
	}

	if err := validateOptionsDuplicates(options); err != nil {
		return err
	}

	return nil
}

func validateOptionsCount(options []option.Option) error {
	if len(options) < MinOptionsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MinOptionsCount)
	}

	if len(options) > MaxOptionsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MaxOptionsCount)
	}

	return nil
}

func validateOptionsContainsEmpty(options []option.Option) error {
	for i := range options {
		if options[i].ID() == uuid.Nil {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	return nil
}

func validateOptionsDuplicates(options []option.Option) error {
	for i := range options {
		for j := i + 1; j < len(options); j++ {
			if options[i].ID() == options[j].ID() {
				return fmt.Errorf("%w: invalid value", ErrInvalid)
			}
		}
	}

	return nil
}
