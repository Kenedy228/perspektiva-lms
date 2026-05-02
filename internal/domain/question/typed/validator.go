package typed

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question/typed/blank"
)

func validatePlaceholders(placeholders []string, blanks []blank.Blank) error {
	if err := validatePlaceholdersCount(placeholders); err != nil {
		return err
	}

	if err := validatePlaceholdersAndBlanksCountMatch(placeholders, blanks); err != nil {
		return err
	}

	if err := validatePlaceholdersAndBlanksPlaceholdersMatch(placeholders, blanks); err != nil {
		return err
	}

	return nil
}

func validatePlaceholdersCount(placeholders []string) error {
	if len(placeholders) < minPlaceholdersCount {
		return fmt.Errorf("%w, детали: число заполнителей в тексте должно быть не менее %d штук", ErrInvalid, minPlaceholdersCount)
	}

	if len(placeholders) > maxPlaceholdersCount {
		return fmt.Errorf("%w, детали: число заполнителей в тексте должно быть не более %d штук", ErrInvalid, maxPlaceholdersCount)
	}

	return nil
}

func validatePlaceholdersAndBlanksCountMatch(placeholders []string, blanks []blank.Blank) error {
	if len(placeholders) != len(blanks) {
		return fmt.Errorf("%w, детали: число заполнителей в тексте не соответствует числу пар 'заполнитель-ответы'", ErrInvalid)
	}

	return nil
}

func validatePlaceholdersAndBlanksPlaceholdersMatch(placeholders []string, blanks []blank.Blank) error {
	blanksPlaceholders := make(map[string]struct{}, len(blanks))

	for i := range blanks {
		blanksPlaceholders[blanks[i].Placeholder()] = struct{}{}
	}

	for i := range placeholders {
		if _, ok := blanksPlaceholders[placeholders[i]]; !ok {
			return fmt.Errorf("%w, детали: пары 'заполнитель-ответы' не содержат заполнитель из текста", ErrInvalid)
		}
	}

	return nil
}
