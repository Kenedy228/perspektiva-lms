package variant

import "strings"

func normalizeValue(value string) string {
	trimmed := strings.TrimSpace(value)
	return strings.Join(strings.Fields(trimmed), " ")
}
