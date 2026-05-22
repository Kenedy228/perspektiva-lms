package answer

import (
	"fmt"
	"unicode/utf8"

	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
)

const (
	MaxBlanksCount      int = 20
	VariantCharsLimit   int = 1000
	PlaceholderCharsMax int = 255
)

func validateBlanks(blanks []AnswerBlank) error {
	if len(blanks) > MaxBlanksCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MaxBlanksCount)
	}

	seen := make(map[string]struct{}, len(blanks))
	for i := range blanks {
		if err := validateBlank(blanks[i]); err != nil {
			return err
		}

		if _, ok := seen[blanks[i].Placeholder]; ok {
			return fmt.Errorf("%w: invalid value", ErrInvalid)
		}
		seen[blanks[i].Placeholder] = struct{}{}
	}

	return nil
}

func validateBlank(b AnswerBlank) error {
	if !blank.SinglePlaceholderRegexp.MatchString(b.Placeholder) {
		return fmt.Errorf("%w: invalid value (%q)", ErrInvalid, b.Placeholder)
	}

	if utf8.RuneCountInString(b.Placeholder) > PlaceholderCharsMax {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, PlaceholderCharsMax)
	}

	if utf8.RuneCountInString(b.Variant) > VariantCharsLimit {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, VariantCharsLimit)
	}

	return nil
}
