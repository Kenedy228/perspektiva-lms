package sequence

import (
	"errors"
	"gitflic.ru/lms/internal/domain/content"
)

var (
	ErrEmptyItems     = errors.New("empty items")
	ErrNotEnoughItems = errors.New("not enough items")
	ErrTooManyItems   = errors.New("too many items")
	ErrItemDuplicate  = errors.New("duplicate item found")
)

func validateItems(items []content.RichContent) error {
	if len(items) == 0 {
		return ErrEmptyItems
	}

	if len(items) < minItems {
		return ErrNotEnoughItems
	}

	if len(items) > maxItems {
		return ErrTooManyItems
	}

	visited := make(map[content.RichContent]struct{}, len(items))

	for i := range items {
		if _, ok := visited[items[i]]; ok {
			return ErrItemDuplicate
		}

		visited[items[i]] = struct{}{}
	}

	return nil
}
