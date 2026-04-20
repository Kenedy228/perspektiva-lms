package sequence

import (
	"errors"

	"gitflic.ru/lms/internal/domain/question/option"
)

var (
	ErrEmptyItems     = errors.New("empty items")
	ErrNotEnoughItems = errors.New("not enough items")
	ErrTooManyItems   = errors.New("too many items")
	ErrItemDuplicate  = errors.New("duplicate item found")
)

func validateItems(items []option.ContentOption) error {
	if len(items) == 0 {
		return ErrEmptyItems
	}

	if len(items) < minItems {
		return ErrNotEnoughItems
	}

	if len(items) > maxItems {
		return ErrTooManyItems
	}

	visitedOptions := make(map[option.ContentOption]struct{}, len(items))
	for i := range items {
		if _, ok := visitedOptions[items[i]]; ok {
			return ErrItemDuplicate
		}

		visitedOptions[items[i]] = struct{}{}
	}

	return nil
}
