package login

import "regexp"

const (
	minValueCharsCount int = 4
	maxValueCharsCount int = 30
)

var (
	validValuePattern = regexp.MustCompile(`^[a-z0-9._-]+$`)
)
