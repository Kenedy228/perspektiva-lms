package sequence

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question/sequence/option"
)

func validateOptions(options []option.Option) error {
	if err := validateOptionsCount(options); err != nil {
		return err
	}

	return nil
}

func validateOptionsCount(options []option.Option) error {
	if len(options) < minOptionsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать минимум %d опций последовательности", ErrInvalid, minOptionsCount)
	}

	if len(options) > maxOptionsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать максимум %d опций последовательности", ErrInvalid, maxOptionsCount)
	}

	return nil
}
