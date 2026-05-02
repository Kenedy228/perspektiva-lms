package blank

import "regexp"

var (
	singlePlaceholderRegexp = regexp.MustCompile(`^\{\{\S+?\}\}$`)
)

const (
	minVariantsCount        int = 1
	maxVariantsCount        int = 20
	variantAsTextCharsLimit int = 1000
)
