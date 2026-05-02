package name

import (
	"strings"
	"unicode"
)

func getInitial(part string) string {
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
