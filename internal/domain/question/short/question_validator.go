package short

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
)

func validateVariants(variants []question.Content) error {
	if len(variants) < minVariants {
		return fmt.Errorf("%w, детали: количество вариантов ответов должно быть не менее %d штук", ErrInvalidVariants, minVariants)
	}

	if len(variants) > maxVariants {
		return fmt.Errorf("%w, детали: количество вариантов ответов должно быть не более %d штук", ErrInvalidVariants, maxVariants)
	}

	for i := range variants {
		if !variants[i].IsText() {
			return fmt.Errorf("%w, детали: вариант ответа должен быть текстового типа", ErrInvalidVariants)
		}

		for j := i + 1; j < len(variants); j++ {
			if variants[i].Equal(variants[j]) {
				return fmt.Errorf("%w, детали: варианты ответов не должны содержать дубликаты", ErrInvalidVariants)
			}
		}
	}

	return nil
}
