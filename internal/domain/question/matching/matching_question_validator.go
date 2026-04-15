package matching

import (
	"errors"
	"strings"

	"gitflic.ru/lms/internal/domain/content"
	"github.com/google/uuid"
)

var (
	ErrNilOption       = errors.New("option cannot be nil")
	ErrEmptyPrompt     = errors.New("prompt cannot be empty")
	ErrNotEnoughPairs  = errors.New("not enough pairs")
	ErrEmptyPairs      = errors.New("empty pairs")
	ErrTooManyPairs    = errors.New("too many pairs")
	ErrPromptDuplicate = errors.New("found duplicated prompt")
	ErrOptionDuplicate = errors.New("found duplicated option")
	ErrExtraPrompts    = errors.New("found extra prompts")
)

func validatePairs(pairs map[string]content.RichContent, pairsCount int) error {
	if len(pairs) == 0 {
		return ErrEmptyPairs
	}

	if len(pairs) < minPairs {
		return ErrNotEnoughPairs
	}

	if len(pairs) > maxPairs {
		return ErrTooManyPairs
	}

	if len(pairs) != pairsCount {
		if len(pairs) < pairsCount {
			return ErrPromptDuplicate
		}

		if len(pairs) > pairsCount {
			return ErrExtraPrompts
		}
	}

	visitedOptions := make(map[content.RichContent]struct{}, len(pairs))
	for _, option := range pairs {
		if _, ok := visitedOptions[option]; ok {
			return ErrOptionDuplicate
		}

		visitedOptions[option] = struct{}{}
	}

	return nil
}

func validatePairPrompt(prompt string) error {
	if strings.TrimSpace(prompt) == "" {
		return ErrEmptyPrompt
	}

	return nil
}

func validatePairOption(optionID uuid.UUID) error {
	if optionID == uuid.Nil {
		return ErrNilOption
	}

	return nil
}
