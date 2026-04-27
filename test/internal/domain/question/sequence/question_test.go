package sequence_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question/sequence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const maxElements = 200

func TestNewQuestion(t *testing.T) {
	t.Run("less than minElements", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withElementAsText("element").
			build(t, sequence.ErrInvalidElements)
	})

	t.Run("more than maxElements", func(t *testing.T) {
		//Arrange-Assert
		b := newQuestionBuilder()

		for i := range maxElements + 1 {
			b = b.withElementAsText(fmt.Sprintf("%d", i))
		}

		b.build(t, sequence.ErrInvalidElements)
	})

	t.Run("with same content", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withElementAsText("same").
			withElementAsText("same").
			build(t, sequence.ErrInvalidElements)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withText("text").
			withImage(uuid.New()).
			withElementAsText("foo").
			withElementAsText("bar").
			build(t, nil))

		//Assert
		assert.NotEqual(t, q.ID(), uuid.Nil)
		assert.Equal(t, q.Text(), "text")
		assert.NotEqual(t, q.ImageID(), uuid.Nil)
		assert.True(t, q.HasImage())
		assert.Equal(t, len(q.Elements()), 2)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})
}

func TestUpdateElements(t *testing.T) {
	t.Run("update with no err should change items and update updatedAt", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withElementAsText("first").
			withElementAsText("second").
			build(t, nil))

		//Act
		err := q.UpdateElements(mockElements())

		//Assert
		assert.Nil(t, err)
		assert.Equal(t, len(q.Elements()), len(mockElements()))
		assert.NotSame(t, &q.Elements()[0], &mockElements()[0])
		assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
	})

	t.Run("update with err should not change items and not update updatedAt", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withElementAsText("first").
			withElementAsText("second").
			build(t, nil))

		//Act
		err := q.UpdateElements([]sequence.Element{newElementBuilder().
			withContentAsText("incorrect").
			buildNoTest()})

		//Assert
		assert.ErrorIs(t, err, sequence.ErrInvalidElements)
		assert.NotEqual(t, q.Elements(), mockElements())
		assert.Equal(t, q.UpdatedAt(), q.CreatedAt())
	})
}

func TestCheckAnswers(t *testing.T) {
	q := castQuestion(t, newQuestionBuilder().withElementAsText("item1").
		withElementAsText("item2").
		build(t, nil))

	require.Equal(t, q.Elements()[0].Content().Value(), "item1")
	require.Equal(t, q.Elements()[1].Content().Value(), "item2")

	first := q.Elements()[0].ID()
	second := q.Elements()[1].ID()

	t.Run("different length", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withRandomID().
			withRandomID().
			withRandomID().
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("equal len, contains duplicates of valid answer", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(first).
			withID(first).
			withRandomID().
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("diff order, valid answer", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(second).
			withID(first).
			withRandomID().
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("correct", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(first).
			withID(second).
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.True(t, check)
	})
}

func TestCloneQuestion(t *testing.T) {
	//Arrange
	q := castQuestion(t, newQuestionBuilder().withElementAsText("text").
		withElementAsText("another").
		build(t, nil))
	clone := castQuestion(t, q.Clone())

	//Assert
	assert.NotSame(t, &q, &clone)
	assert.NotSame(t, &q.Elements()[0], &clone.Elements()[0])
}
