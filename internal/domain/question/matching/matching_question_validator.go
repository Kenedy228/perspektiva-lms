package matching

import (
	"errors"
	"strings"

	"gitflic.ru/lms/internal/domain/question/option"
)

var (
	ErrEmptyPairs      = errors.New("empty pairs")
	ErrNotEnoughPairs  = errors.New("not enough pairs")
	ErrTooManyPairs    = errors.New("too many pairs")
	ErrDuplicatePrompt = errors.New("found duplicated prompt")
	ErrDuplicateOption = errors.New("found duplicated option")
	ErrEmptyPrompt     = errors.New("prompt cannot be empty")
)

func validatePairs(pairs []PairParams) error {
	if len(pairs) == 0 {
		return ErrEmptyPairs
	}

	if len(pairs) < minPairs {
		return ErrNotEnoughPairs
	}

	if len(pairs) > maxPairs {
		return ErrTooManyPairs
	}

	visitedPrompts := make(map[string]struct{}, len(pairs))
	visitedOptions := make(map[option.ContentOption]struct{}, len(pairs))

	for i := range pairs {
		if _, ok := visitedPrompts[pairs[i].Prompt]; ok {
			return ErrDuplicatePrompt
		}

		if _, ok := visitedOptions[pairs[i].ContentOption]; ok {
			return ErrDuplicateOption
		}

		visitedPrompts[pairs[i].Prompt] = struct{}{}
		visitedOptions[pairs[i].ContentOption] = struct{}{}
	}

	return nil
}

func validatePrompt(prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		return ErrEmptyPrompt
	}

	return nil
}
