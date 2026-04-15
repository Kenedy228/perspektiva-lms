package sequence

import (
	"errors"
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/content"
)

func TestValidateItems(t *testing.T) {
	type fooItem struct {
		cType content.ContentType
		value string
	}

	tests := []struct {
		name  string
		items []fooItem
		err   error
	}{
		{
			name:  "empty items",
			items: []fooItem{},
			err:   ErrEmptyItems,
		},
		{
			name:  "not enough items",
			items: []fooItem{{cType: content.ContentTypeImage, value: "url"}},
			err:   ErrNotEnoughItems,
		},
		{
			name:  "with duplicates",
			items: []fooItem{{cType: content.ContentTypeImage, value: "url"}, {cType: content.ContentTypeImage, value: "url"}},
			err:   ErrItemDuplicate,
		},
		{
			name:  "without duplicates",
			items: []fooItem{{cType: content.ContentTypeText, value: "url"}, {cType: content.ContentTypeImage, value: "url"}},
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]content.RichContent, 0, len(tt.items))

			for i := range tt.items {
				item, err := content.New(tt.items[i].cType, tt.items[i].value)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				items = append(items, item)
			}

			err := validateItems(items)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}

func TestValidateItemsWithExceededLimit(t *testing.T) {
	tests := []struct {
		name string
		size int
		err  error
	}{
		{
			name: "size 19",
			size: 19,
			err:  nil,
		},
		{
			name: "size 20",
			size: 20,
			err:  nil,
		},
		{
			name: "size 21",
			size: 21,
			err:  ErrTooManyItems,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]content.RichContent, 0, tt.size)			

			for i := range tt.size {
				item, err := content.New(content.ContentTypeImage, fmt.Sprintf("%d", i))

				if err != nil {
					t.Errorf("expected err nil, got %v", err)
				}

				items = append(items, item)
			}

			err := validateItems(items)

			if !errors.Is(err, tt.err) {
				t.Errorf("expected err %v, got %v", tt.err, err)
			}
		})
	}
}
