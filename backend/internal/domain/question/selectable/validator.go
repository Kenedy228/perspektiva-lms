package selectable

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
)

func validateBase(b *base.Base) error {
	if b == nil {
		return fmt.Errorf("%w: база вопроса должна существовать", ErrInvalid)
	}

	return nil
}

func validateOptions(options []option.Option) error {
	if err := validateOptionsCount(options); err != nil {
		return err
	}

	if err := validateOptionsDuplicates(options); err != nil {
		return err
	}

	if err := validateOptionsContainsEmpty(options); err != nil {
		return err
	}

	if err := validateCorrectOptionsCount(options); err != nil {
		return err
	}

	return nil
}

func validateOptionsCount(options []option.Option) error {
	count := len(options)

	if count < MinOptionsCount {
		return fmt.Errorf("%w: вопрос с выбором должен содержать минимум %d опции (текущее количество - %d)", ErrInvalid, MinOptionsCount, count)
	}

	if count > MaxOptionsCount {
		return fmt.Errorf("%w: вопрос с выбором должен содержать максимум %d опций (текущее количество - %d)", ErrInvalid, MaxOptionsCount, count)
	}

	return nil
}

func validateCorrectOptionsCount(options []option.Option) error {
	count := countCorrectOptions(options)

	if count < MinCorrectOptionsCount {
		return fmt.Errorf("%w: вопрос с выбором должен содержать минимум %d правильную опцию", ErrInvalid, MinCorrectOptionsCount)
	}

	return nil
}

func validateOptionsDuplicates(options []option.Option) error {
	for i := range options {
		for j := i + 1; j < len(options); j++ {
			if options[i].ID() == options[j].ID() {
				return fmt.Errorf("%w: опция %s добавлена в вопрос больше одного раза", ErrInvalid, options[i].ID())
			}
		}
	}

	return nil
}

func validateOptionsContainsEmpty(options []option.Option) error {
	for i := range options {
		if options[i].IsZero() {
			return fmt.Errorf("%w: вопрос с выбором не должен содержать пустых опций", ErrInvalid)
		}
	}

	return nil
}
