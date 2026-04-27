package short_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/short"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const maxVariants = 200

func TestNewQuestion(t *testing.T) {
	t.Run("zero length", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().build(t, short.ErrInvalidVariants)
	})

	t.Run("more than maxVariants", func(t *testing.T) {
		//Arrange-Assert
		b := newQuestionBuilder()

		for i := range maxVariants + 1 {
			b = b.withVariant(fmt.Sprintf("%d", i))
		}

		b.build(t, short.ErrInvalidVariants)
	})

	t.Run("with same content", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().withVariant("same").
			withVariant("same").
			build(t, short.ErrInvalidVariants)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withText("text").
			withImage(uuid.New()).
			withVariant("foo").
			withVariant("bar").
			build(t, nil))

		//Assert
		assert.NotEqual(t, q.ID(), uuid.Nil)
		assert.Equal(t, q.Text(), "text")
		assert.NotEqual(t, q.ImageID(), uuid.Nil)
		assert.True(t, q.HasImage())
		assert.Equal(t, len(q.Variants()), 2)
		assert.Equal(t, q.CreatedAt(), q.UpdatedAt())
	})
}

func TestUpdateItems(t *testing.T) {
	t.Run("update with no err should change items and update updatedAt", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withVariant("first").
			withVariant("second").
			build(t, nil))

		//Act
		err := q.UpdateVariants(mockVariants())

		//Assert
		assert.Nil(t, err)
		assert.Equal(t, len(q.Variants()), len(mockVariants()))
		assert.NotSame(t, &q.Variants()[0], &mockVariants()[0])
		assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
	})

	t.Run("update with err should not change items and not update updatedAt", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().withVariant("first").
			withVariant("second").
			build(t, nil))

		//Act
		err := q.UpdateVariants([]question.Content{})

		//Assert
		assert.ErrorIs(t, err, short.ErrInvalidVariants)
		assert.NotEqual(t, q.Variants(), mockVariants())
		assert.Equal(t, q.UpdatedAt(), q.CreatedAt())
	})
}

func TestCheckAnswers(t *testing.T) {
	q := castQuestion(t, newQuestionBuilder().withVariant("item1").
		withVariant("item2").
		build(t, nil))

	t.Run("q contains answer", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withInput("item1").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.True(t, check)
	})

	t.Run("q contains answer in diff case", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withInput("ITEM1").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.True(t, check)
	})

	t.Run("q doesn't contain answer", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withInput("item3").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("answer is empty", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withInput("").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})
}

func TestCloneQuestion(t *testing.T) {
	//Arrange
	q := castQuestion(t, newQuestionBuilder().withVariant("text").
		withVariant("another").
		build(t, nil))

	//Act
	clone := castQuestion(t, q.Clone())

	//Assert
	assert.NotSame(t, &q, &clone)
	assert.NotSame(t, &q.Variants()[0], &clone.Variants()[0])
}
