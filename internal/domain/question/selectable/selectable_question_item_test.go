package selectable

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/option"
	"github.com/stretchr/testify/assert"
)

func TestNewItem(t *testing.T) {
	c := makeContent("c")

	t.Run("testing correct item", func(t *testing.T) {
		item := NewItem(makeItemParam("c", true))

		assert.Equal(t, item.Content(), c)
		assert.True(t, item.IsCorrect())
	})

	t.Run("testing incorrect item", func(t *testing.T) {
		item := NewItem(makeItemParam("c", false))

		assert.Equal(t, item.Content(), c)
		assert.False(t, item.IsCorrect())
	})
}

func TestItemEqual(t *testing.T) {
	baseItem := NewItem(ItemParams{
		Content:   makeContent("content"),
		IsCorrect: true,
	})

	tests := []struct {
		name    string
		content string
		correct bool
		equal   bool
	}{
		{
			name:    "different contents",
			content: "diff",
			correct: true,
			equal:   false,
		},
		{
			name:    "correct false, same content",
			content: "content",
			correct: false,
			equal:   true,
		},
		{
			name:    "correct true, same content",
			content: "content",
			correct: true,
			equal:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewItem(makeItemParam(tt.content, tt.correct))

			if tt.equal {
				assert.True(t, item.Equal(baseItem))
			} else {
				assert.False(t, item.Equal(baseItem))
			}
		})
	}
}

func makeContent(s string) option.ContentOption {
	c, _ := option.NewContentOption(option.ContentTypeText, s)
	return c
}

func makeItemParam(s string, correct bool) ItemParams {
	return ItemParams{
		Content:   makeContent(s),
		IsCorrect: correct,
	}
}
