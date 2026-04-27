package typed

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
)

func validatePlaceholder(placeholder string) error {
	if !singlePlaceholderRegexp.MatchString(placeholder) {
		return fmt.Errorf("%w, детали: некорректный формат заполнителя", ErrInvalidPlaceholder)
	}

	return nil
}

func validateVariants(variants []question.Content) error {
	if len(variants) < minVariants {
		return fmt.Errorf("%w, детали: количество ответов для заполнителя должно быть не менее %d штук", ErrInvalidVariants, minVariants)
	}

	if len(variants) > maxVariants {
		return fmt.Errorf("%w, детали: количество ответов для заполнителя должно быть не более %d штук", ErrInvalidVariants, maxVariants)
	}

	for i := range variants {
		if !variants[i].IsText() {
			return fmt.Errorf("%w, детали: ответы для заполнителя должны быть в текстовом формате", ErrInvalidVariants)
		}
	}

	return nil
}
