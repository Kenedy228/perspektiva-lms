package sequence

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/content"
	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockRichContent(id int) content.RichContent {
	p := content.Params{
		Type:  content.ContentTypeImage,
		Value: fmt.Sprintf("%d", id),
	}
	rc, _ := content.New(p)
	return rc
}

func TestNewItem(t *testing.T) {
	c := mockRichContent(1)
	item, err := NewItem(c)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, item.ID())
	assert.Equal(t, c, item.Content())
}

func TestNew(t *testing.T) {
	// Генерируем слишком большой массив контента
	tooManyItems := make([]content.RichContent, maxItems+1)
	for i := range tooManyItems {
		tooManyItems[i] = mockRichContent(i)
	}

	// Валидный массив из 3 элементов
	validItems := []content.RichContent{
		mockRichContent(1),
		mockRichContent(2),
		mockRichContent(3),
	}

	tests := []struct {
		name    string
		params  *Params
		wantErr error
	}{
		{
			name: "success valid sequence",
			params: &Params{
				Text:  "Расставьте события по порядку",
				Items: validItems,
			},
			wantErr: nil,
		},
		{
			name: "error empty items",
			params: &Params{
				Text:  "Вопрос без вариантов",
				Items: []content.RichContent{},
			},
			wantErr: ErrEmptyItems,
		},
		{
			name: "error not enough items",
			params: &Params{
				Text:  "Слишком мало вариантов",
				Items: []content.RichContent{mockRichContent(1)}, // Только 1 элемент
			},
			wantErr: ErrNotEnoughItems,
		},
		{
			name: "error too many items",
			params: &Params{
				Text:  "Слишком много вариантов",
				Items: tooManyItems,
			},
			wantErr: ErrTooManyItems,
		},
		{
			name: "error duplicate items",
			params: &Params{
				Text: "Одинаковые варианты",
				Items: []content.RichContent{
					mockRichContent(1),
					mockRichContent(1), // Дубликат
				},
			},
			wantErr: ErrItemDuplicate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q, err := New(tt.params)

			assert.ErrorIs(t, err, tt.wantErr)

			if tt.wantErr == nil {
				require.NotNil(t, q)
				sequenceQ, ok := q.(*SequenceQuestion)
				require.True(t, ok)

				assert.Equal(t, question.TypeSequence, sequenceQ.Type())
				assert.Len(t, sequenceQ.Items(), len(tt.params.Items))
			}
		})
	}
}

func TestSequenceQuestion_UpdateItems(t *testing.T) {
	initialItems := []content.RichContent{
		mockRichContent(1),
		mockRichContent(2),
	}

	params := &Params{
		Text:  "Базовый текст",
		Items: initialItems,
	}
	q, err := New(params)
	require.NoError(t, err)

	sequenceQ, ok := q.(*SequenceQuestion)
	require.True(t, ok)

	t.Run("success update items", func(t *testing.T) {
		newItems := []content.RichContent{
			mockRichContent(3),
			mockRichContent(4),
			mockRichContent(5),
		}

		err := sequenceQ.UpdateItems(newItems)

		assert.NoError(t, err)
		assert.Len(t, sequenceQ.Items(), 3)
		assert.Equal(t, newItems[0], sequenceQ.Items()[0].Content())
	})

	t.Run("error update leaves state untouched", func(t *testing.T) {
		oldItems := sequenceQ.Items()

		// Пытаемся обновить невалидными данными (менее 2 элементов)
		invalidItems := []content.RichContent{mockRichContent(6)}
		err := sequenceQ.UpdateItems(invalidItems)

		assert.ErrorIs(t, err, ErrNotEnoughItems)

		// Убеждаемся, что старые элементы не затерлись
		assert.Equal(t, oldItems, sequenceQ.Items())
	})
}

func TestSequenceQuestion_Encapsulation(t *testing.T) {
	initialItems := []content.RichContent{
		mockRichContent(1),
		mockRichContent(2),
	}

	params := &Params{
		Text:  "Текст",
		Items: initialItems,
	}

	q, err := New(params)
	require.NoError(t, err)
	sequenceQ := q.(*SequenceQuestion)

	t.Run("getter clone", func(t *testing.T) {
		gotItems := sequenceQ.Items()
		require.Len(t, gotItems, 2)

		// Пытаемся затереть первый элемент в полученном срезе
		gotItems[0] = Item{}

		// Проверяем, что внутреннее состояние не изменилось (оригинал не пострадал)
		assert.NotEqual(t, uuid.Nil, sequenceQ.Items()[0].ID())
	})
}
