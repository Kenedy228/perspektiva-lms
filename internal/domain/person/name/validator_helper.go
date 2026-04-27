package name

import (
	"fmt"
	"unicode"
)

func validatePartStructure(field, value string) error {
	runes := []rune(value)

	if len(runes) == 0 {
		return fmt.Errorf("%w, детали: %s не может быть пустым", ErrInvalid, field)
	}

	if len(runes) > maxPartLen {
		return fmt.Errorf("%w, детали: количество символов в %s превышает лимит %d символов", ErrInvalid, field, maxPartLen)
	}

	seenLetter := false
	prevWasSpace := false
	prevWasSeparator := false

	for i, r := range runes {
		switch {
		case isCyrillicLetter(r):
			seenLetter = true
			prevWasSpace = false
			prevWasSeparator = false
		case isSeparatorLetter(r):
			if i == 0 || i == len(runes)-1 {
				return fmt.Errorf("%w, детали: %s содержит разделитель в начале или в конце", ErrInvalid, field)
			}

			if prevWasSeparator || prevWasSpace {
				return fmt.Errorf("%w, детали: %s содержит некорректную последовательность разделителей", ErrInvalid, field)
			}

			prevWasSeparator = true
			prevWasSpace = false
		case isSpaceLetter(r):
			if prevWasSeparator {
				return fmt.Errorf("%w, детали: %s содержит некорректную последовательность символов", ErrInvalid, field)
			}

			prevWasSpace = true
			prevWasSeparator = false
		default:
			if unicode.IsDigit(r) || unicode.IsNumber(r) {
				return fmt.Errorf("%w, детали: %s содержит цифры", ErrInvalid, field)
			}

			return fmt.Errorf("%w, детали: %s содержит недопустимый символ", ErrInvalid, field)
		}
	}

	if !seenLetter {
		return fmt.Errorf("%w, детали: %s должно содержать хотя бы одну букву", ErrInvalid, field)
	}

	if prevWasSpace || prevWasSeparator {
		return fmt.Errorf("%w, детали: %s содержит некорректное окончание", ErrInvalid, field)
	}

	return nil
}

func isCyrillicLetter(r rune) bool {
	return unicode.IsLetter(r) && unicode.Is(unicode.Cyrillic, r)
}

func isSeparatorLetter(r rune) bool {
	return r == '-' || r == '\''
}

func isSpaceLetter(r rune) bool {
	return r == ' '
}
