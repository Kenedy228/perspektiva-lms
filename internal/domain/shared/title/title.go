package title

import (
	"errors"
	"strings"
)

type Title string

var (
	ErrEmptyTitle = errors.New("empty title")
)

func NewTitle(s string) (Title, error) {
	if strings.TrimSpace(s) == "" {
		return Title(""), ErrEmptyTitle
	}

	return Title(s), nil
}

func (t Title) String() string {
	return string(t)
}
