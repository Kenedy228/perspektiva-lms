package education

import "strings"

func normalizeValue(education string) string {
	trimmed := strings.TrimSpace(education)
	return strings.Join(strings.Fields(trimmed), " ")
}
