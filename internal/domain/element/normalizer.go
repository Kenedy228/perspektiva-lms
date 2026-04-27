package element

import "strings"

func normalizeTitle(s string) string {
	trimmed := strings.TrimSpace(s)
	fields := strings.Fields(trimmed)
	joined := strings.Join(fields, " ")
	return joined
}
