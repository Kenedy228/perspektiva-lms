package account

import "strings"

func normalize(s string) string {
	trimmed := strings.TrimSpace(s)
	joined := strings.Join(strings.Fields(trimmed), "")
	return strings.ToLower(joined)
}
