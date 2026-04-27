package typed

import "regexp"

const (
	minPlaceholders = 2
	maxPlaceholders = 20
	minVariants     = 1
	maxVariants     = 20
)

var (
	singlePlaceholderRegexp = regexp.MustCompile(`^\{\{\S+?\}\}$`)
	inTextPlaceholderRegexp = regexp.MustCompile(`\{\{\S+?\}\}`)
)
