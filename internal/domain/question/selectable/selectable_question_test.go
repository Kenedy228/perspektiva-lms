package selectable

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelectableQuestion(t *testing.T) {
	tooManyItems := make([]ItemParams, 0, maxItems+1)
	for i := 0; i <= maxItems; i++ {
		p := makeItemParam(fmt.Sprintf("option %d", i), true)
		tooManyItems = append(tooManyItems, p)
	}

	tests := []struct {
		name   string
		params Params
		err    error
	}{
		{
			name: "success valid question",
			params: Params{
				Text: makeText("Выберите цвета светофора"),
				Items: []ItemParams{
					makeItemParam("Красный", true),
					makeItemParam("Зеленый", true),
					makeItemParam("Фиолетовый", false),
				},
			},
			err: nil,
		},
		{
			name: "error empty items",
			params: Params{
				Text:  makeText("Вопрос без вариантов"),
				Items: []ItemParams{},
			},
			err: ErrEmptyItems,
		},
		{
			name: "error not enough items",
			params: Params{
				Text: makeText("Один вариант ответа"),
				Items: []ItemParams{
					makeItemParam("Одинокий ответ", true),
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
			name: "error no correct items",
			params: Params{
				Text: makeText("Где правильный?"),
				Items: []ItemParams{
					makeItemParam("Неправильно", false),
					makeItemParam("Тоже нет", false),
				},
			},
			err: ErrNoCorrectItem,
		},
		{
			name: "error duplicated content",
			params: Params{
				Text: makeText("Дубликат контента?"),
				Items: []ItemParams{
					makeItemParam("Неправильно", false),
					makeItemParam("Неправильно", true),
				},
			},
			err: ErrDuplicateItem,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := NewSelectableQuestion(tt.params)

			assert.ErrorIs(t, err, tt.err)

			if tt.err == nil {
				require.NotNil(t, q)
				selectableQ, ok := q.(*SelectableQuestion)
				require.True(t, ok)

				assert.Equal(t, question.TypeSelectable, selectableQ.Type())
				assert.Len(t, selectableQ.Items(), len(tt.params.Items))
				for i := range selectableQ.Items() {
					assert.True(t, selectableQ.Items()[i].Equal(mapItems(tt.params.Items)[i]))
				}
			}
		})
	}
}

func TestSelectableQuestion_UpdateItems(t *testing.T) {
	params := Params{
		Text: makeText("Изначальный текст"),
		Items: []ItemParams{
			makeItemParam("Вариант 1", true),
			makeItemParam("Вариант 2", false),
		},
	}
	q, err := NewSelectableQuestion(params)
	require.NoError(t, err)

	selectableQ, ok := q.(*SelectableQuestion)
	require.True(t, ok)

	t.Run("success update items", func(t *testing.T) {
		newItems := []ItemParams{
			makeItemParam("Новый 1", true),
			makeItemParam("Новый 2", true),
			makeItemParam("Новый 3", false),
		}

		err := selectableQ.UpdateItems(newItems)

		assert.NoError(t, err)
		assert.Len(t, selectableQ.Items(), 3)
		for i := range selectableQ.Items() {
			assert.True(t, selectableQ.Items()[i].Equal(mapItems(newItems)[i]))
		}
	})

	t.Run("error update leaves state untouched", func(t *testing.T) {
		oldItems := selectableQ.Items()

		invalidItems := []ItemParams{
			makeItemParam("Ошибка 1", false),
			makeItemParam("Ошибка 2", false),
		}
		err := selectableQ.UpdateItems(invalidItems)

		assert.ErrorIs(t, err, ErrNoCorrectItem)

		for i := range selectableQ.Items() {
			assert.True(t, selectableQ.Items()[i].Equal(oldItems[i]))
		}
	})
}

func TestHasItem(t *testing.T) {
	params := Params{
		Text: makeText("Изначальный текст"),
		Items: []ItemParams{
			makeItemParam("Вариант 1", true),
			makeItemParam("Вариант 2", false),
		},
	}
	q, err := NewSelectableQuestion(params)
	require.NoError(t, err)

	selectableQ, ok := q.(*SelectableQuestion)
	require.True(t, ok)

	t.Run("has item", func(t *testing.T) {
		item := NewItem(makeItemParam("Вариант 1", false))

		assert.True(t, selectableQ.HasItem(item))
	})

	t.Run("has not item", func(t *testing.T) {
		item := NewItem(makeItemParam("Вариант 3", false))

		assert.False(t, selectableQ.HasItem(item))
	})
}

func makeText(s string) question.QText {
	text, _ := question.NewQText(s)
	return text
}
