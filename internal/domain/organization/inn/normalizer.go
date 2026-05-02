package inn

import "strings"

func normalizeCode(code string) string {
	return strings.ReplaceAll(code, " ", "")
}
