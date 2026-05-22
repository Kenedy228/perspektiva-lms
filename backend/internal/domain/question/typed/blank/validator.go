package blank

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
)

func validatePlaceholder(placeholder string) error {
	if !SinglePlaceholderRegexp.MatchString(placeholder) {
		return fmt.Errorf("%w: invalid value (%q)", ErrInvalid, placeholder)
	}
	return nil
}

func validateVariants(variants []text.Text) error {
	if err := validateVariantsCount(variants); err != nil {
		return err
	}
	if err := validateVariantsContainsEmpty(variants); err != nil {
		return err
	}
	if err := validateVariantsDuplicates(variants); err != nil {
		return err
	}
	if err := validateVariantsCharsLimit(variants); err != nil {
		return err
	}
	return nil
}

func validateVariantsCount(variants []text.Text) error {
	if len(variants) < MinVariantsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MinVariantsCount)
	}
	if len(variants) > MaxVariantsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MaxVariantsCount)
	}
	return nil
}

func validateVariantsContainsEmpty(variants []text.Text) error {
	for i := range variants {
		if len(variants[i].Value()) == 0 {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}
	return nil
}

func validateVariantsDuplicates(variants []text.Text) error {
	for i := range variants {
		for j := i + 1; j < len(variants); j++ {
			v1 := strings.TrimSpace(strings.ToLower(variants[i].Value()))
			v2 := strings.TrimSpace(strings.ToLower(variants[j].Value()))
			if v1 == v2 {
				return fmt.Errorf("%w: invalid value", ErrInvalid)
			}
		}
	}
	return nil
}

func validateVariantsCharsLimit(variants []text.Text) error {
	for i := range variants {
		if utf8.RuneCountInString(variants[i].Value()) > VariantCharsLimit {
			return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, VariantCharsLimit)
		}
	}
	return nil
}
