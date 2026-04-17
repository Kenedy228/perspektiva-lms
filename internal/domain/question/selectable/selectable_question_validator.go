package selectable

import (
	"errors"

	"gitflic.ru/lms/internal/domain/question/option"
)

var (
	ErrEmptyItems            = errors.New("empty items")
	ErrNotEnoughItems        = errors.New("not enough items")
	ErrTooManyItems          = errors.New("too many items")
	ErrNoCorrectItem         = errors.New("no correct item provided")
	ErrNotEnoughCorrectItems = errors.New("not enough correct items provided")
	ErrDuplicateItem         = errors.New("duplicated item")
)

func validateItems(items []ItemParams) error {
	if len(items) == 0 {
		return ErrEmptyItems
	}

	if len(items) < minItems {
		return ErrNotEnoughItems
	}

	if len(items) > maxItems {
		return ErrTooManyItems
	}

	correctCount := 0
	visitedContent := make(map[option.ContentOption]struct{}, len(items))
	for i := range items {
		if _, ok := visitedContent[items[i].Content]; ok {
			return ErrDuplicateItem
		}

		visitedContent[items[i].Content] = struct{}{}

		if items[i].IsCorrect {
			correctCount++
		}
	}

	if correctCount == 0 {
		return ErrNoCorrectItem
	}

	if correctCount < minCorrectItems {
		return ErrNotEnoughCorrectItems
	}

	return nil
}
