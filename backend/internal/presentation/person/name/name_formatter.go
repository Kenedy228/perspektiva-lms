package name

import (
	"fmt"
	"strings"
	"unicode"

	"gitflic.ru/lms/backend/internal/domain/person/name"
)

type FullNameFormatter struct{}

func (f FullNameFormatter) Format(n name.Name) string {
	if n.MiddleName() == "" {
		return fmt.Sprintf("%s %s", n.LastName(), n.FirstName())
	}

	return fmt.Sprintf("%s %s %s", n.LastName(), n.FirstName(), n.MiddleName())
}

type WithInitialsFormatter struct{}

// TODO: доделать форматтер этот

func (f WithInitialsFormatter) Format(n name.Name) string {
	getInitial := func(part string) string {
		if part == "" {
			return ""
		}

		var builder strings.Builder
		runes := []rune(part)

		builder.WriteRune(unicode.ToUpper(runes[0]))
		builder.WriteRune('.')

		for i := 1; i < len(runes); i++ {
			if runes[i] == '-' && i+1 < len(runes) {
				builder.WriteRune('-')
				builder.WriteRune(unicode.ToUpper(runes[i+1]))
				builder.WriteRune('.')
			}
			if runes[i] == ' ' && i+1 < len(runes) {
				builder.WriteRune(unicode.ToUpper(runes[i+1]))
				builder.WriteRune('.')
			}
		}

		return builder.String()
	}

	if n.MiddleName() == "" {
		return fmt.Sprintf("%s %s", n.LastName(), getInitial(n.FirstName()))
	}

	return fmt.Sprintf("%s %s%s", n.LastName(), getInitial(n.FirstName()), getInitial(n.MiddleName()))
}
