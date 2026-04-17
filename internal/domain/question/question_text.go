package question

import (
	"errors"
	"strings"
)

type QText string

var (
	ErrEmptyText = errors.New("empty question text")
)

func NewQText(s string) (QText, error) {
	if strings.TrimSpace(s) == "" {
		return QText(""), ErrEmptyText
	}

	return QText(s), nil
}

func (t QText) String() string {
	return string(t)
}

func (t QText) Equal(other QText) bool {
	return t.String() == other.String()
}
