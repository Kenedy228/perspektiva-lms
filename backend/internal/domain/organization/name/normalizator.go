package name

import "strings"

func normalizeValue(value string) string {
	return strings.Join(strings.Fields(value), " ")
}
