package blank

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/internal/domain/question/content"
)

func validatePlaceholder(placeholder string) error {
	if err := validatePlaceholderFormat(placeholder); err != nil {
		return err
	}

	return nil
}

func validatePlaceholderFormat(placeholder string) error {
	if !singlePlaceholderRegexp.MatchString(placeholder) {
		return fmt.Errorf("%w, детали: некорректный формат заполнителя", ErrInvalid)
	}

	return nil
}

func validateVariants(variants []content.Content) error {
	if err := validateVariantsCount(variants); err != nil {
		return err
	}

	if err := validateVariantsAllowedFormat(variants); err != nil {
		return err
	}

	if err := validateVariantsAsTextCharsLimit(variants); err != nil {
		return err
	}

	return nil
}

func validateVariantsCount(variants []content.Content) error {
	if len(variants) < minVariantsCount {
		return fmt.Errorf("%w, детали: количество ответов для заполнителя должно быть не менее %d штук", ErrInvalid, minVariantsCount)
	}

	if len(variants) > maxVariantsCount {
		return fmt.Errorf("%w, детали: количество ответов для заполнителя должно быть не более %d штук", ErrInvalid, maxVariantsCount)
	}

	return nil
}

func validateVariantsAllowedFormat(variants []content.Content) error {
	for i := range variants {
		if !variants[i].IsText() {
			return fmt.Errorf("%w, детали: вариант ответа должен быть в текстовом формате", ErrInvalid)
		}
	}

	return nil
}

func validateVariantsAsTextCharsLimit(variants []content.Content) error {
	for i := range variants {
		if utf8.RuneCountInString(variants[i].Value()) > variantAsTextCharsLimit {
			return fmt.Errorf("%w, детали: вариант в текстовом формате должен содержать не более %d символов", ErrInvalid, variantAsTextCharsLimit)
		}
	}

	return nil
}
