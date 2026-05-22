package login

import "regexp"

const (
	MinValueCharsCount int = 3
	MaxValueCharsCount int = 10
)

var (
	validValuePattern = regexp.MustCompile(`^[A-Za-z0-9-]+$`)
)
