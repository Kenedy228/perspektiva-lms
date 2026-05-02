package selectable

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question/selectable/option"
)

func validateOptions(options []option.Option) error {
	if err := validateOptionsCount(options); err != nil {
		return err
	}

	if err := validateCorrectOptionsCount(options); err != nil {
		return err
	}

	return nil
}

func validateOptionsCount(options []option.Option) error {
	if len(options) < minOptionsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать минимум %d вариантов ответа", ErrInvalid, minOptionsCount)
	}

	if len(options) > maxOptionsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать максимум %d вариантов ответа", ErrInvalid, maxOptionsCount)
	}

	return nil
}

func validateCorrectOptionsCount(options []option.Option) error {
	count := countCorrectOptions(options)

	if count < minCorrectOptionsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать минимум %d правильных вариантов ответа", ErrInvalid, minCorrectOptionsCount)
	}

	return nil
}
