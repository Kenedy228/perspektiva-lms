package sequence

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tooManyItems := make([]option.ContentOption, 0, maxItems+1)
	for i := range maxItems + 1 {
		tooManyItems = append(tooManyItems, makeItem(fmt.Sprintf("%d", i)))
	}

	tests := []struct {
		name   string
		params Params
		err    error
	}{
		{
			name: "success valid sequence",
			params: Params{
				Text: makeText("Расставьте события по порядку"),
				Items: []option.ContentOption{
					makeItem("1"),
					makeItem("2"),
					makeItem("3"),
				},
			},
			err: nil,
		},
		{
			name: "error empty items",
			params: Params{
				Text:  makeText("Вопрос без вариантов"),
				Items: []option.ContentOption{},
			},
			err: ErrEmptyItems,
		},
		{
			name: "error not enough items",
			params: Params{
				Text: makeText("Слишком мало вариантов"),
				Items: []option.ContentOption{
					makeItem("1"),
				},
			},
			err: ErrNotEnoughItems,
		},
		{
			name: "error too many items",
			params: Params{
				Text:  makeText("Слишком много вариантов"),
				Items: tooManyItems,
			},
			err: ErrTooManyItems,
		},
		{
			name: "error duplicate items",
			params: Params{
				Text: "Одинаковые варианты",
				Items: []option.ContentOption{
					makeItem("1"),
					makeItem("1"),
				},
			},
			err: ErrItemDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, q)
				sequenceQ, ok := q.(*SequenceQuestion)
				require.True(t, ok)

				assert.Equal(t, question.TypeSequence, sequenceQ.Type())
				assert.Len(t, sequenceQ.Items(), len(tt.params.Items))
				assert.Equal(t, sequenceQ.Items(), tt.params.Items)
			}
		})
	}
}

func TestSequenceQuestion_UpdateItems(t *testing.T) {
	initialItems := []option.ContentOption{
		makeItem("1"),
		makeItem("2"),
	}

	params := Params{
		Text:  makeText("Базовый текст"),
		Items: initialItems,
	}
	q, err := New(params)
	require.NoError(t, err)

	sequenceQ, ok := q.(*SequenceQuestion)
	require.True(t, ok)

	t.Run("success update items", func(t *testing.T) {
		newItems := []option.ContentOption{
			makeItem("3"),
			makeItem("4"),
			makeItem("5"),
		}

		err := sequenceQ.UpdateItems(newItems)

		assert.NoError(t, err)
		assert.Len(t, sequenceQ.Items(), 3)
		assert.Equal(t, sequenceQ.Items(), newItems)
	})

	t.Run("error update leaves state untouched", func(t *testing.T) {
		oldItems := sequenceQ.Items()

		invalidItems := []option.ContentOption{makeItem("6")}
		err := sequenceQ.UpdateItems(invalidItems)

		assert.ErrorIs(t, err, ErrNotEnoughItems)
		assert.Equal(t, sequenceQ.Items(), oldItems)
	})
}

func TestHasItem(t *testing.T) {
	initialItems := []option.ContentOption{
		makeItem("1"),
		makeItem("2"),
	}

	params := Params{
		Text:  makeText("Базовый текст"),
		Items: initialItems,
	}
	q, err := New(params)
	require.NoError(t, err)

	sequenceQ, ok := q.(*SequenceQuestion)
	require.True(t, ok)

	t.Run("has item", func(t *testing.T) {
		assert.True(t, sequenceQ.HasItem(makeItem("1")))
	})

	t.Run("has not item", func(t *testing.T) {
		assert.False(t, sequenceQ.HasItem(makeItem("10")))
	})
}

func makeText(s string) question.QText {
	text, _ := question.NewQText(s)
	return text
}

func makeItem(s string) option.ContentOption {
	item, _ := option.NewContentOption(option.ContentTypeText, s)
	return item
}
