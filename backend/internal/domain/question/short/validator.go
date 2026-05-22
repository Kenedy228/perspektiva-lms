package short

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
)

func validateVariants(variants []variant.Variant) error {
	if err := validateVariantsCount(variants); err != nil {
		return err
	}

	if err := validateVariantsContainsNonInitialized(variants); err != nil {
		return err
	}

	return nil
}

func validateVariantsCount(variants []variant.Variant) error {
	count := len(variants)

	if count < MinVariantsCount {
		return fmt.Errorf("%w: для создания вопроса необходимо минимум %d вариантов ответов, текущее количество - %d", ErrInvalid, MinVariantsCount, count)
	}

	if count > MaxVariantsCount {
		return fmt.Errorf("%w: для создания вопроса необходимо максимум %d вариантов ответов, текущее количество - %d", ErrInvalid, MaxVariantsCount, count)
	}

	return nil
}

func validateVariantsContainsNonInitialized(variants []variant.Variant) error {
	for i := range variants {
		if variants[i].IsZero() {
			return fmt.Errorf("%w: вопрос не может содержать пустых вариантов ответов", ErrInvalid)
		}
	}

	return nil
}
