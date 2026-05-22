package typed

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
)

func uniquePlaceholdersInText(text string) []string {
	matches := InTextPlaceholderRegexp.FindAllString(text, -1)
	set := make(map[string]struct{}, len(matches))
	var unique []string

	for _, m := range matches {
		if _, exists := set[m]; !exists {
			set[m] = struct{}{}
			unique = append(unique, m)
		}
	}
	return unique
}

func validatePlaceholdersAndBlanks(placeholders []string, blanks []blank.Blank) error {
	if err := validatePlaceholdersCount(placeholders); err != nil {
		return err
	}
	if err := validateBlanksUniqueness(blanks); err != nil {
		return err
	}
	if err := validateMatch(placeholders, blanks); err != nil {
		return err
	}
	return nil
}

func validatePlaceholdersCount(placeholders []string) error {
	if len(placeholders) < MinPlaceholdersCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MinPlaceholdersCount)
	}
	if len(placeholders) > MaxPlaceholdersCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, MaxPlaceholdersCount)
	}
	return nil
}

func validateBlanksUniqueness(blanks []blank.Blank) error {
	seen := make(map[string]struct{}, len(blanks))
	for i := range blanks {
		p := blanks[i].Placeholder()
		if _, ok := seen[p]; ok {
			return fmt.Errorf("%w: invalid value (%q)", ErrInvalid, p)
		}
		seen[p] = struct{}{}
	}
	return nil
}

func validateMatch(placeholders []string, blanks []blank.Blank) error {
	if len(placeholders) != len(blanks) {
		return fmt.Errorf("%w: invalid value (%d) (%d)", ErrInvalid, len(placeholders), len(blanks))
	}

	textSet := make(map[string]struct{}, len(placeholders))
	for _, p := range placeholders {
		textSet[p] = struct{}{}
	}

	for _, b := range blanks {
		if _, ok := textSet[b.Placeholder()]; !ok {
			return fmt.Errorf("%w: invalid value (%q)", ErrInvalid, b.Placeholder())
		}
	}

	return nil
}
