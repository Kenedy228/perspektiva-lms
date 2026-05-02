package short

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question/short/variant"
)

func validateVariants(variants []variant.Variant) error {
	if err := validateVariantsCount(variants); err != nil {
		return err
	}

	return nil
}

func validateVariantsCount(variants []variant.Variant) error {
	if len(variants) < minVariantsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать не менее %d вариантов ответа", ErrInvalid, minVariantsCount)
	}

	if len(variants) > maxVariantsCount {
		return fmt.Errorf("%w, детали: вопрос должен содержать не более %d вариантов ответов", ErrInvalid, maxVariantsCount)
	}

	return nil
}
