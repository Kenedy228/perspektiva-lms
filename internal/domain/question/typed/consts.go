package typed

import "regexp"

const (
	minPlaceholdersCount int = 2
	maxPlaceholdersCount int = 20
)

var (
	inTextPlaceholderRegexp = regexp.MustCompile(`\{\{\S+?\}\}`)
)
