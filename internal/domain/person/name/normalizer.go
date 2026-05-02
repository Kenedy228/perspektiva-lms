package name

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.Russian)

func normalize(s string) string {
	trimmed := strings.TrimSpace(s)
	cleaned := strings.Join(strings.Fields(trimmed), " ")

	if cleaned == "" {
		return ""
	}

	parts := strings.Split(cleaned, "-")

	for i := range parts {
		parts[i] = titleCaser.String(parts[i])
	}

	return strings.Join(parts, "-")
}
