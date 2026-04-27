package block_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/block"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustNewBlockWithNElems(t *testing.T, n int) *block.Block {
	t.Helper()

	b := newBlockBuilder().
		withTitle("Блок").
		build(t, nil)

	for i := 0; i < n; i++ {
		require.NoError(t, b.InsertElementAt(i, newElement()))
	}

	return b
}

func TestBlock_InsertElementAt_Success(t *testing.T) {
	b := newBlockBuilder().
		withTitle("Блок").
		withElement().
		withElement().
		build(t, nil)
	require.NotNil(t, b)

	el := newElement()

	// вставка в середину
	err := b.InsertElementAt(1, el)
	require.NoError(t, err)

	elems := b.Elements()
	require.Len(t, elems, 3)
	assert.Equal(t, el.ID(), elems[1].ID())
}

func TestBlock_InsertElementAt_InvalidPosition(t *testing.T) {
	b := mustNewBlockWithNElems(t, 2)

	tests := []struct {
		name string
		pos  int
	}{
		{"negative", -1},
		{"too big", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.InsertElementAt(tt.pos, newElement())
			assert.ErrorIs(t, err, block.ErrInvalid)
		})
	}
}

func TestBlock_InsertElementAt_NilElement(t *testing.T) {
	b := mustNewBlockWithNElems(t, 1)

	err := b.InsertElementAt(1, nil)
	assert.ErrorIs(t, err, block.ErrInvalid)
}

func TestBlock_InsertElementAt_DuplicatesForbidden(t *testing.T) {
	b := mustNewBlockWithNElems(t, 1)
	existing := b.Elements()[0]

	err := b.InsertElementAt(1, existing)
	assert.ErrorIs(t, err, block.ErrInvalid)
}

func TestBlock_RemoveElementAt_Success(t *testing.T) {
	b := mustNewBlockWithNElems(t, 3)

	err := b.RemoveElementAt(1)
	require.NoError(t, err)

	elems := b.Elements()
	require.Len(t, elems, 2)
}

func TestBlock_RemoveElementAt_InvalidPosition(t *testing.T) {
	b := mustNewBlockWithNElems(t, 2)

	tests := []struct {
		name string
		pos  int
	}{
		{"negative", -1},
		{"equal len", 2},
		{"greater len", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.RemoveElementAt(tt.pos)
			assert.ErrorIs(t, err, block.ErrInvalid)
		})
	}
}

func TestBlock_UpdateElement_Success(t *testing.T) {
	b := mustNewBlockWithNElems(t, 2)
	elems := b.Elements()
	old := elems[0]

	updated := old.Clone()
	updated.ChangeTitle("new title")

	err := b.UpdateElement(updated)
	require.NoError(t, err)

	assert.Equal(t, b.Elements()[0].Title(), "new title")
}

func TestBlock_UpdateElement_NotFound(t *testing.T) {
	b := mustNewBlockWithNElems(t, 1)

	updated := newElement() // с другим ID

	err := b.UpdateElement(updated)
	assert.ErrorIs(t, err, block.ErrInvalid)
}

func TestBlock_MoveFromTo_Success(t *testing.T) {
	b := newBlockBuilder().
		withTitle("Блок").
		withElement(). // 0
		withElement(). // 1
		withElement(). // 2
		build(t, nil)
	require.NotNil(t, b)

	elems := b.Elements()
	id0 := elems[0].ID()
	id1 := elems[1].ID()
	id2 := elems[2].ID()

	t.Run("move middle to end", func(t *testing.T) {
		err := b.MoveFromTo(1, 2)
		require.NoError(t, err)

		e := b.Elements()
		require.Len(t, e, 3)
		assert.Equal(t, id0, e[0].ID())
		assert.Equal(t, id2, e[1].ID())
		assert.Equal(t, id1, e[2].ID())
	})

	t.Run("move first to last", func(t *testing.T) {
		err := b.MoveFromTo(0, 2)
		require.NoError(t, err)

		e := b.Elements()
		require.Len(t, e, 3)
		assert.Equal(t, id0, e[2].ID())
		assert.Equal(t, id1, e[1].ID())
		assert.Equal(t, id2, e[0].ID())
	})
}

func TestBlock_MoveFromTo_InvalidPositions(t *testing.T) {
	b := mustNewBlockWithNElems(t, 2)

	tests := []struct {
		name string
		from int
		to   int
	}{
		{"from negative", -1, 0},
		{"from >= len", 2, 0},
		{"to negative", 0, -1},
		{"to >= len", 0, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := b.MoveFromTo(tt.from, tt.to)
			assert.ErrorIs(t, err, block.ErrInvalid)
		})
	}
}

func TestBlock_Clone(t *testing.T) {
	// given
	b := newBlockBuilder().
		withTitle("Исходный блок").
		withElement().
		withElement().
		build(t, nil)
	require.NotNil(t, b)

	// when
	clone := b.Clone()

	// then: базовые свойства
	require.NotNil(t, clone)
	assert.Equal(t, b.ID(), clone.ID(), "Clone должен сохранять id")
	assert.Equal(t, b.Title(), clone.Title())
	assert.Len(t, clone.Elements(), len(b.Elements()))

	// then: элементы — глубокая копия
	origElems := b.Elements()
	cloneElems := clone.Elements()

	for i := range origElems {
		assert.Equal(t, origElems[i].ID(), cloneElems[i].ID(),
			"Clone: элементы должны иметь те же ID")
		assert.NotSame(t, origElems[i], cloneElems[i],
			"Clone: элементы должны быть разными указателями (глубокая копия)")
	}

	// модифицируем клон и убеждаемся, что оригинал не меняется
	cloneElems[0].ChangeTitle("другое") // или любой доступный мутатор Element
	assert.NotEqual(t, cloneElems[0].Title(), origElems[0].Title(),
		"изменения в клоне не должны затрагивать оригинал")
}

func TestBlock_Copy(t *testing.T) {
	// given
	b := newBlockBuilder().
		withTitle("Исходный блок").
		withElement().
		withElement().
		build(t, nil)
	require.NotNil(t, b)

	// when
	copyBlock, err := b.Copy()
	require.NoError(t, err)

	// then: базовые свойства
	require.NotNil(t, copyBlock)
	assert.NotEqual(t, b.ID(), copyBlock.ID(), "Copy должен генерировать новый id")
	assert.Equal(t, b.Title(), copyBlock.Title())
	assert.Len(t, copyBlock.Elements(), len(b.Elements()))

	// then: элементы — глубокая копия
	origElems := b.Elements()
	copyElems := copyBlock.Elements()

	for i := range origElems {
		assert.Equal(t, origElems[i].ID(), copyElems[i].ID(),
			"Copy: элементы должны иметь те же ID")
		assert.NotSame(t, origElems[i], copyElems[i],
			"Copy: элементы должны быть разными указателями (глубокая копия)")
	}

	// модифицируем копию и убеждаемся, что оригинал не меняется
	copyElems[0].ChangeTitle("другое") // или любой мутатор
	assert.NotEqual(t, copyElems[0].Title(), origElems[0].Title(),
		"изменения в копии не должны затрагивать оригинал")
}
