package version

import "strings"

func normalizeTitle(title string) string {
	title = strings.TrimSpace(title)
	return strings.Join(strings.Fields(title), " ")
}
