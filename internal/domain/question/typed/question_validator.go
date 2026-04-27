package typed

import (
	"fmt"
)

func validatePlaceholders(text string, blanks []Blank) error {
	foundPlaceholders := inTextPlaceholderRegexp.FindAllString(text, -1)

	if len(foundPlaceholders) < minPlaceholders {
		return fmt.Errorf("%w, детали: число заполнителей в тексте должно быть не менее %d штук", ErrInvalidText, minPlaceholders)
	}

	if len(foundPlaceholders) > maxPlaceholders {
		return fmt.Errorf("%w, детали: число заполнителей в тексте должно быть не более %d штук", ErrInvalidText, maxPlaceholders)
	}

	for i := range foundPlaceholders {
		for j := i + 1; j < len(foundPlaceholders); j++ {
			if foundPlaceholders[i] == foundPlaceholders[j] {
				return fmt.Errorf("%w, детали: заполнители не должны содержать дубликатов", ErrInvalidText)
			}
		}
	}

	if len(foundPlaceholders) != len(blanks) {
		return fmt.Errorf("%w, детали: число заполнителей в тексте не соответствует числу пар 'заполнитель-ответы'", ErrInvalidBlanks)
	}

	for i := range foundPlaceholders {
		count := 0
		for j := range blanks {
			if foundPlaceholders[i] == blanks[j].Placeholder() {
				count++
			}
		}

		if count != 1 {
			return fmt.Errorf("%w, детали: пары 'заполнитель-ответы' содержат дубликаты заполнителей из текста", ErrInvalidBlanks)
		}
	}

	return nil
}
