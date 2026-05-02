package base_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/base"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	b, err := base.New(makeTitle("title"))

	//Assert
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, b.ID())
	assert.Equal(t, "title", b.Title().Value())
	assert.False(t, b.HasAttachment())
}

func TestChangeTitle(t *testing.T) {
	//Arrange
	b, err := base.New(makeTitle("title"))
	require.NoError(t, err)

	//Act
	newTitle := makeTitle("new")
	b.ChangeTitle(newTitle)

	//Assert
	assert.Equal(t, b.Title(), newTitle)
}

func TestChangeAttachment(t *testing.T) {
	//Arrange
	b, err := base.New(makeTitle("title"))
	require.NoError(t, err)

	//Act
	a := makeAttachment("attachment")
	b.ChangeAttachment(a)
	attachment, ok := b.Attachment()

	//Assert
	assert.True(t, ok)
	assert.Equal(t, attachment, a)
	assert.True(t, b.HasAttachment())
}

func TestRemoveAttachment(t *testing.T) {
	//Arrange
	b, err := base.New(makeTitle("title"))
	require.NoError(t, err)

	//Act
	b.RemoveAttachment()
	_, ok := b.Attachment()

	//Assert
	assert.False(t, ok)
	assert.False(t, b.HasAttachment())
}

func TestClone(t *testing.T) {
	t.Run("с вложением", func(t *testing.T) {
		//Arrange
		b, err := base.New(makeTitle("title"))
		require.NoError(t, err)

		//Act
		b.ChangeAttachment(makeAttachment("attachment"))
		clone := b.Clone()

		originalAttachment, ok := b.Attachment()
		require.True(t, ok)
		cloneAttachment, ok := clone.Attachment()
		require.True(t, ok)

		//Assert
		assert.NotSame(t, b, clone)
		assert.Equal(t, b.ID(), clone.ID())
		assert.Equal(t, b.Title(), clone.Title())
		assert.True(t, b.HasAttachment())
		assert.Equal(t, originalAttachment, cloneAttachment)
	})

	t.Run("без вложения", func(t *testing.T) {
		//Arrange
		b, err := base.New(makeTitle("title"))
		require.NoError(t, err)

		//Act
		clone := b.Clone()

		//Assert
		assert.NotSame(t, b, clone)
		assert.Equal(t, b.ID(), clone.ID())
		assert.Equal(t, b.Title(), clone.Title())
		assert.False(t, b.HasAttachment())
		assert.Equal(t, b.HasAttachment(), clone.HasAttachment())
	})
}
