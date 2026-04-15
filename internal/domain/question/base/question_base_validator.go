package base

import (
	"errors"
	"strings"
)

var (
	ErrEmptyText        = errors.New("text cannot be empty")
	ErrEmptyDescription = errors.New("description cannot be empty")
)

func validateText(text string) error {
	if strings.TrimSpace(text) == "" {
		return ErrEmptyText
	}

	return nil
}

func validateDescription(description string) error {
	if strings.TrimSpace(description) == "" {
		return ErrEmptyDescription
	}

	return nil
}
