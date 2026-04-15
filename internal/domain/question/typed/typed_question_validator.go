package typed

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyMark                = errors.New("empty mark")
	ErrMarkDuplicate            = errors.New("duplicate mark in text")
	ErrEmptyBlankAnswer         = errors.New("empty blank answer")
	ErrNoBlankAnswers           = errors.New("no answers for blank")
	ErrNoPlaceholders           = errors.New("no placeholders in text")
	ErrPlaceholderCountMismatch = errors.New("text and blanks contains different count of placeholders")
	ErrPlaceholderMissing       = errors.New("placeholder missing")
	ErrTooManyPlaceholders      = errors.New("too many placeholders")
)

func validateBlank(mark string, answers []string) error {
	if strings.TrimSpace(mark) == "" {
		return ErrEmptyMark
	}

	if len(answers) == 0 {
		return ErrNoBlankAnswers
	}

	for i := range answers {
		if strings.TrimSpace(answers[i]) == "" {
			return ErrEmptyBlankAnswer
		}
	}

	return nil
}

func validatePlaceholders(text string, placeholdersCount int, blanks map[string][]string) error {
	if placeholdersCount == 0 {
		return ErrNoPlaceholders
	}

	if placeholdersCount > maxPlaceholders {
		return ErrTooManyPlaceholders
	}

	if len(blanks) != placeholdersCount {
		return ErrPlaceholderCountMismatch
	}

	for mark := range blanks {
		placeholder := fmt.Sprintf("[%s]", mark)
		entriesCount := strings.Count(text, placeholder)

		if entriesCount == 0 {
			return ErrPlaceholderMissing
		}

		if entriesCount > 1 {
			return ErrMarkDuplicate
		}
	}

	return nil
}
