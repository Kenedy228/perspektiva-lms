package short

import (
	"fmt"
	"strings"

	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
)

func validateVariants(variants []variant.Variant) error {
	if err := validateVariantsCount(variants); err != nil {
		return err
	}

	if err := validateVariantsContainsEmpty(variants); err != nil {
		return err
	}

	if err := validateVariantsDuplicates(variants); err != nil {
		return err
	}

	return nil
}

func validateVariantsCount(variants []variant.Variant) error {
	if len(variants) < MinVariantsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MinVariantsCount)
	}

	if len(variants) > MaxVariantsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MaxVariantsCount)
	}

	return nil
}

func validateVariantsContainsEmpty(variants []variant.Variant) error {
	for i := range variants {
		if variants[i].IsZero() {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
	}

	return nil
}

func validateVariantsDuplicates(variants []variant.Variant) error {
	for i := range variants {
		for j := i + 1; j < len(variants); j++ {
			v1 := strings.TrimSpace(strings.ToLower(variants[i].Text().Value()))
			v2 := strings.TrimSpace(strings.ToLower(variants[j].Text().Value()))

			if v1 == v2 {
				return fmt.Errorf("%w: invalid value", ErrInvalid)
			}
		}
	}

	return nil
}
