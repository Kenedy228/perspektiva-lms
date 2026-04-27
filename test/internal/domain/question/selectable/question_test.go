package selectable_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question/selectable"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const maxOptions = 200

func TestNewQuestion(t *testing.T) {
	t.Run("less than minAnswers", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withOptionAsText("item", true).
			build(t, selectable.ErrInvalidOptions)
	})

	t.Run("more than maxAnswers", func(t *testing.T) {
		//Arrange-Assert
		b := newQuestionBuilder()

		for i := range maxOptions + 1 {
			b = b.withOptionAsText(fmt.Sprintf("%d", i), true)
		}

		b.build(t, selectable.ErrInvalidOptions)
	})

	t.Run("with same content", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withOptionAsText("same", true).
			withOptionAsText("same", false).
			build(t, selectable.ErrInvalidOptions)
	})

	t.Run("less than minCorrect", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withOptionAsText("foo", false).
			withOptionAsText("bar", false).
			build(t, selectable.ErrInvalidOptions)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withText("text").
			withImage(uuid.New()).
			withOptionAsText("foo", true).
			withOptionAsText("bar", false).
			build(t, nil))

		//Assert
		assert.NotEqual(t, q.ID(), uuid.Nil)
		assert.Equal(t, q.Text(), "text")
		assert.NotEqual(t, q.ImageID(), uuid.Nil)
		assert.True(t, q.HasImage())
		assert.Equal(t, len(q.Options()), 2)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})
}

func TestUpdateOptions(t *testing.T) {
	t.Run("update with no err should change items and update updatedAt", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withOptionAsText("first", true).
			withOptionAsText("second", false).
			build(t, nil))

		//Act
		err := q.UpdateOptions(mockOptions())

		//Assert
		assert.Nil(t, err)
		assert.Equal(t, len(q.Options()), len(mockOptions()))
		assert.NotSame(t, &q.Options()[0], &mockOptions()[0])
		assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
	})

	t.Run("update with err should not change items and not update updatedAt", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withOptionAsText("first", true).
			withOptionAsText("second", false).
			build(t, nil))

		//Act
		err := q.UpdateOptions([]selectable.Option{newOptionBuilder().withContentAsText("incorrect").
			withCorrect(false).
			buildNoTest()})

		//Assert
		assert.ErrorIs(t, err, selectable.ErrInvalidOptions)
		assert.NotEqual(t, q.Options(), mockOptions())
		assert.Equal(t, q.UpdatedAt(), q.CreatedAt())
	})
}

func TestCheckAnswers(t *testing.T) {
	q := castQuestion(t, newQuestionBuilder().withOptionAsText("item1", true).
		withOptionAsText("item2", false).
		build(t, nil))

	require.True(t, q.Options()[0].IsCorrect())
	require.False(t, q.Options()[1].IsCorrect())

	truthy := q.Options()[0].ID()
	falsy := q.Options()[1].ID()

	t.Run("contains id, that q doesn't have", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withRandomID().
			withRandomID().
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer contains falsy id", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(falsy).
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer contains falsy id", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(falsy).
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer contains falsy and truthy id", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(falsy).
			withID(truthy).
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer contains duplicate truthy id", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(truthy).
			withID(truthy).
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer contains only truthy id", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withID(truthy).
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.True(t, check)
	})
}

func TestCloneQuestion(t *testing.T) {
	//Arrange
	q := castQuestion(t, newQuestionBuilder().withOptionAsText("text", true).
		withOptionAsText("another", false).
		build(t, nil))
	clone := castQuestion(t, q.Clone())

	assert.Equal(t, q.ID(), clone.ID())
	assert.NotSame(t, &q, &clone)
	assert.NotSame(t, &q.Options()[0], &clone.Options()[0])
	assert.Equal(t, q.UpdatedAt(), clone.UpdatedAt())
	assert.Equal(t, q.CreatedAt(), clone.CreatedAt())
}
