package sequence

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
)

func validateBase(b *base.Base) error {
	if err := validateBaseRequired(b); err != nil {
		return err
	}

	return nil
}

func validateBaseRequired(b *base.Base) error {
	if b == nil {
		return fmt.Errorf("%w: база вопроса должна существовать", ErrInvalid)
	}

	return nil
}

func validateOptions(options []option.Option) error {
	if err := validateOptionsCount(options); err != nil {
		return err
	}

	if err := validateOptionsContainsEmpty(options); err != nil {
		return err
	}

	return nil
}

func validateOptionsCount(options []option.Option) error {
	count := len(options)

	if count < MinOptionsCount {
		return fmt.Errorf("%w: вопрос должен содержать минимум %d опций (текущее количество опций - %d)", ErrInvalid, MinOptionsCount, count)
	}

	if count > MaxOptionsCount {
		return fmt.Errorf("%w: вопрос должен содержать максимум %d опций (текущее количество опций - %d)", ErrInvalid, MaxOptionsCount, count)
	}

	return nil
}

func validateOptionsContainsEmpty(options []option.Option) error {
	for i := range options {
		if options[i].IsZero() {
			return fmt.Errorf("%w: вопрос не должен содержать пустых опций", ErrInvalid)
		}
	}

	return nil
}
