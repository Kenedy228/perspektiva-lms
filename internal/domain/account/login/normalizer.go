package login

import "strings"

func normalizeValue(value string) string {
	trimmed := strings.TrimSpace(value)
	return strings.ToLower(strings.ReplaceAll(trimmed, " ", ""))
}
