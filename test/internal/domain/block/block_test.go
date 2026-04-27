package block_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlockNew_Success(t *testing.T) {
	t.Run("valid title and no elements", func(t *testing.T) {
		b := newBlockBuilder().
			withTitle("Мой блок").
			build(t, nil)

		require.NotNil(t, b)
		assert.Equal(t, "Мой блок", b.Title())
		assert.Len(t, b.Elements(), 0)
	})

	t.Run("valid title and some elements", func(t *testing.T) {
		b := newBlockBuilder().
			withTitle("Блок с элементами").
			withElement().
			withElement().
			build(t, nil)

		require.NotNil(t, b)
		assert.Equal(t, "Блок с элементами", b.Title())
		assert.Len(t, b.Elements(), 2)
	})
}

func TestBlockNew_InvalidTitle(t *testing.T) {
	tests := []struct {
		name  string
		title string
	}{
		{"empty string", ""},
		{"spaces only", "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := newBlockBuilder().
				withTitle(tt.title).
				build(t, block.ErrInvalid)

			assert.Nil(t, b)
		})
	}
}

func TestBlockNew_TitleTooLong(t *testing.T) {
	// Соберём строку длиной titleCharsLimit+1.
	long := make([]rune, 1e5) // если константа экспортируемая
	for i := range long {
		long[i] = 'a'
	}

	b := newBlockBuilder().
		withTitle(string(long)).
		build(t, block.ErrInvalid)

	assert.Nil(t, b)
}

func TestBlockNew_ElementsLimitExceeded(t *testing.T) {
	builder := newBlockBuilder().withTitle("Блок с лишними элементами")

	for i := 0; i < 1e5; i++ { // если константа экспортируется
		builder = builder.withElement()
	}

	b := builder.build(t, block.ErrInvalid)
	assert.Nil(t, b)
}

func TestBlockNew_ElementsContainNil(t *testing.T) {
	builder := newBlockBuilder().
		withTitle("Блок с nil элементом").
		withElement()

	// вручную добавляем nil
	builder.elements = append(builder.elements, nil)

	b := builder.build(t, block.ErrInvalid)
	assert.Nil(t, b)
}

func TestBlockNew_ElementsWithDuplicates(t *testing.T) {
	el := newElement()

	builder := newBlockBuilder().
		withTitle("Блок с дубликатами").
		withElement()

	// подменим первый реальный элемент на заранее созданный
	builder.elements[0] = el
	// добавим тот же указатель ещё раз
	builder.elements = append(builder.elements, el)

	b := builder.build(t, block.ErrInvalid)
	assert.Nil(t, b)
}
