//go:build legacy
// +build legacy

package block_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlock_FullLifecycle(t *testing.T) {
	// 1. Создание блока
	b := newBlockBuilder().
		withTitle("Блок").
		build(t, nil)
	require.NotNil(t, b)
	assert.Len(t, b.Elements(), 0)

	// 2. Добавление элементов
	el1 := newElement()
	el2 := newElement()
	el3 := newElement()

	require.NoError(t, b.InsertElementAt(0, el1))
	require.NoError(t, b.InsertElementAt(1, el2))
	require.NoError(t, b.InsertElementAt(2, el3))

	elems := b.Elements()
	require.Len(t, elems, 3)
	assert.Equal(t, el1.ID(), elems[0].ID())
	assert.Equal(t, el2.ID(), elems[1].ID())
	assert.Equal(t, el3.ID(), elems[2].ID())

	// 3. Перемещение el2 в конец
	require.NoError(t, b.MoveFromTo(1, 2))

	elems = b.Elements()
	require.Len(t, elems, 3)
	assert.Equal(t, el1.ID(), elems[0].ID())
	assert.Equal(t, el3.ID(), elems[1].ID())
	assert.Equal(t, el2.ID(), elems[2].ID())

	// 4. Обновление el3 (аналогично TODO выше — нужен элемент с тем же ID)
	updatedEl3 := el3.Clone()
	updatedEl3.ChangeTitle("new title")
	require.NoError(t, b.UpdateElement(updatedEl3))
	elems = b.Elements()
	assert.Equal(t, elems[1].Title(), "new title")

	// 5. Удаление первого элемента
	require.NoError(t, b.RemoveElementAt(0))

	elems = b.Elements()
	require.Len(t, elems, 2)
	assert.Equal(t, el3.ID(), elems[0].ID())
	assert.Equal(t, el2.ID(), elems[1].ID())
}
