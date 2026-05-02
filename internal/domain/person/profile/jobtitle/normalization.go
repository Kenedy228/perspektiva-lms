package jobtitle

import "strings"

func normalizeJobTitle(jobTitle string) string {
	trimmed := strings.TrimSpace(jobTitle)
	return strings.Join(strings.Fields(trimmed), " ")
}
