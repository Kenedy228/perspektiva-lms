package name

import "strings"

func normalize(s string) string {
	trimmed := strings.TrimSpace(s)
	return strings.Join(strings.Fields(trimmed), " ")
}
