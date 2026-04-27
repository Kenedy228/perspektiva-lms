package typed_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/typed"
	"github.com/stretchr/testify/assert"
)

func TestNewQuestion(t *testing.T) {
	t.Run("less than min placeholders", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().
			withText("hello {{placeholder}} not enough").
			withBlank("placeholder", "value").
			build(t, typed.ErrInvalidText)
	})

	t.Run("more than maxPlaceholders", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().
			withText("{{1}}{{2}}{{3}}{{4}}{{5}}{{6}}{{7}}{{8}}{{9}}{{11}}{{12}}{{13}}{{14}}{{15}}{{16}}{{17}}{{18}}{{19}}{{21}}{{22}}{{23}}").
			withBlank("placeholder", "value").
			build(t, typed.ErrInvalidText)
	})

	t.Run("with placeholder duplicates", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().
			withText("text with {{duplicate}} {{duplicate}}").
			withBlank("duplicate", "value").
			build(t, typed.ErrInvalidText)
	})

	t.Run("with missing placeholder", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().
			withText("text with {{placeholder}} {{missing}}").
			withBlank("placeholder", "value").
			build(t, typed.ErrInvalidBlanks)
	})

	t.Run("with missing placeholder", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().
			withText("text with {{placeholder}} {{missing}}").
			withBlank("placeholder", "value").
			build(t, typed.ErrInvalidBlanks)
	})

	t.Run("blanks with duplicates", func(t *testing.T) {
		//Arrange-Assert
		newQuestionBuilder().
			withText("text with {{placeholder}} {{missing}}").
			withBlank("placeholder", "value").
			withBlank("placeholder", "value").
			build(t, typed.ErrInvalidBlanks)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().
			withText("text with {{placeholder}} {{another}}").
			withBlank("placeholder", "value").
			withBlank("another", "value").
			build(t, nil))

		//Assert
		assert.Equal(t, q.Text(), "text with {{placeholder}} {{another}}")
		assert.Equal(t, len(q.Blanks()), 2)
	})
}

func TestReplaceContent(t *testing.T) {
	t.Run("success update", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().
			withText("text with {{placeholder}} {{another}}").
			withBlank("placeholder", "value").
			withBlank("another", "value").
			build(t, nil))

		blanks := []typed.Blank{
			newBlankBuilder().withPlaceholder("new").
				withVariant("answer").
				buildNoTest(),
			newBlankBuilder().withPlaceholder("new2").
				withVariant("answer").
				buildNoTest(),
		}

		//Act
		err := q.ReplaceContent("text with {{new}} {{new2}}", blanks)

		//Assert
		assert.Nil(t, err)
		assert.Equal(t, q.Text(), "text with {{new}} {{new2}}")
		assert.Equal(t, len(q.Blanks()), 2)
		assert.NotSame(t, &q.Blanks()[0], &blanks[0])
		assert.True(t, q.UpdatedAt().After(q.CreatedAt()))
	})

	t.Run("err update", func(t *testing.T) {
		//Arrange
		q := castQuestion(t, newQuestionBuilder().
			withText("text with {{placeholder}} {{another}}").
			withBlank("placeholder", "value").
			withBlank("another", "value").
			build(t, nil))

		//Act
		err := q.ReplaceContent("text with {{new}} {{new2}}", []typed.Blank{})

		//Assert
		assert.ErrorIs(t, err, typed.ErrInvalidBlanks)
		assert.Equal(t, q.Text(), "text with {{placeholder}} {{another}}")
		assert.Equal(t, len(q.Blanks()), 2)
		assert.Equal(t, q.UpdatedAt(), q.CreatedAt())
	})
}

func TestCheckAnswer(t *testing.T) {
	q := castQuestion(t, newQuestionBuilder().
		withText("text with {{placeholder}} {{another}}").
		withBlank("placeholder", "value").
		withBlank("another", "value").
		build(t, nil))

	t.Run("different len", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withBlank("{{placeholder}}", "value").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("same len, but has different placeholders", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withBlank("{{placeholder}}", "value").
			withBlank("{{diff}}", "value").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("same len, same placeholders, diff values", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withBlank("{{placeholder}}", "value123").
			withBlank("{{another}}", "value").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.False(t, check)
	})

	t.Run("same len, same placeholders, same values", func(t *testing.T) {
		//Arrange
		answer := newAnswerBuilder().withBlank("{{placeholder}}", "value").
			withBlank("{{another}}", "value").
			build()

		//Act
		check := q.CheckAnswer(answer)

		//Assert
		assert.True(t, check)
	})
}

func TestClone(t *testing.T) {
	//Arrange
	q := castQuestion(t, newQuestionBuilder().
		withText("text with {{placeholder}} {{another}}").
		withBlank("placeholder", "value").
		withBlank("another", "value").
		build(t, nil))

	//Act
	clone := castQuestion(t, q.Clone())

	//Assert
	assert.Equal(t, q.ID(), clone.ID())
	assert.Equal(t, q.Text(), clone.Text())
	assert.Equal(t, q.ImageID(), clone.ImageID())
	assert.NotSame(t, &q.Blanks()[0], &clone.Blanks()[0])
	assert.Equal(t, q.CreatedAt(), clone.CreatedAt())
	assert.Equal(t, q.UpdatedAt(), clone.UpdatedAt())
}
