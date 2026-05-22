package login

import "strings"

func normalizeValue(value string) string {
	trimmed := strings.TrimSpace(value)
	return strings.ReplaceAll(trimmed, " ", "")
}
