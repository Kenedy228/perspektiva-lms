//go:build legacy
// +build legacy

package typed_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/typed"
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/question/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQuestion(t *testing.T) {
	t.Run("валидное создание", func(t *testing.T) {
		// Arrange & Act
		q := newQuestionBuilder().
			withTitle("text with {{placeholder}} and {{another}}").
			withBlank("{{placeholder}}", "value1").
			withBlank("{{another}}", "value2").
			build(t)

		// Assert
		assert.Equal(t, "text with {{placeholder}} and {{another}}", q.Title().Value())
		assert.Len(t, q.Blanks(), 2)
	})

	t.Run("меньше необходимого числа плейсхолдеров", func(t *testing.T) {
		// Arrange & Act
		_, err := newQuestionBuilder().
			withTitle("hello {{placeholder}} not enough").
			withBlank("{{placeholder}}", "value").
			buildWithError()

		// Assert
		assert.ErrorIs(t, err, typed.ErrInvalid)
		assert.Contains(t, err.Error(), "не менее 2 штук")
	})

	t.Run("больше необходимого числа плейсхолдеров", func(t *testing.T) {
		// Arrange & Act
		qb := newQuestionBuilder().withTitle(
			"{{1}} {{2}} {{3}} {{4}} {{5}} {{6}} {{7}} {{8}} {{9}} {{10}} " +
				"{{11}} {{12}} {{13}} {{14}} {{15}} {{16}} {{17}} {{18}} {{19}} {{20}} {{21}}",
		)
		for i := 1; i <= 21; i++ {
			qb = qb.withBlank("{{placeholder}}", "value")
		}

		_, err := qb.buildWithError()

		// Assert
		assert.ErrorIs(t, err, typed.ErrInvalid)
		assert.Contains(t, err.Error(), "не более 20 штук")
	})

	t.Run("отсутствует плейсхолдер в бланке", func(t *testing.T) {
		// Arrange & Act
		_, err := newQuestionBuilder().
			withTitle("text with {{placeholder}} {{missing}}").
			withBlank("{{placeholder}}", "value").
			withBlank("{{wrong}}", "value"). // Бланк не совпадает с плейсхолдером в тексте
			buildWithError()

		// Assert
		assert.ErrorIs(t, err, typed.ErrInvalid)
	})

	t.Run("забыли плейсхолдер в бланке", func(t *testing.T) {
		// Arrange & Act
		_, err := newQuestionBuilder().
			withTitle("text with {{placeholder}} {{missing}}").
			withBlank("{{placeholder}}", "value"). // Забыли второй бланк
			buildWithError()

		// Assert
		assert.ErrorIs(t, err, typed.ErrInvalid)
	})
}

func TestReplaceContent(t *testing.T) {
	t.Run("успешное обновление", func(t *testing.T) {
		// Arrange
		q := newQuestionBuilder().
			withTitle("text with {{placeholder}} {{another}}").
			withBlank("{{placeholder}}", "value").
			withBlank("{{another}}", "value").
			build(t)

		newTitleStr := "text with {{new}} {{new2}}"
		newTitle, _ := title.New(makeContent(newTitleStr))

		newBlanks := []blank.Blank{
			makeBlank("{{new}}", "answer"),
			makeBlank("{{new2}}", "answer"),
		}

		// Act
		err := q.ReplaceContent(newTitle, newBlanks)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, "text with {{new}} {{new2}}", q.Title().Value())
		assert.Len(t, q.Blanks(), 2)

		// Проверяем защитное копирование
		assert.NotSame(t, &q.Blanks()[0], &newBlanks[0])
	})

	t.Run("пустые бланки", func(t *testing.T) {
		// Arrange
		q := newQuestionBuilder().
			withTitle("text with {{placeholder}} {{another}}").
			withBlank("{{placeholder}}", "value").
			withBlank("{{another}}", "value").
			build(t)

		newTitle, _ := title.New(makeContent("text with {{new}} {{new2}}"))

		// Act (передаем пустые бланки — это вызовет ошибку)
		err := q.ReplaceContent(newTitle, []blank.Blank{})

		// Assert
		require.ErrorIs(t, err, typed.ErrInvalid)

		// Состояние не должно было измениться!
		assert.Equal(t, "text with {{placeholder}} {{another}}", q.Title().Value())
		assert.Len(t, q.Blanks(), 2)
	})
}

func TestClone(t *testing.T) {
	// Arrange
	q := newQuestionBuilder().
		withTitle("text with {{placeholder}} {{another}}").
		withBlank("{{placeholder}}", "value").
		withBlank("{{another}}", "value").
		build(t)

	// Act
	clone, ok := q.Clone().(*typed.Question)
	require.True(t, ok, "cloned question must be *typed.Question")

	// Assert
	assert.Equal(t, q.ID(), clone.ID())
	assert.Equal(t, q.Title().Value(), clone.Title().Value())
	assert.Equal(t, q.Instruction(), clone.Instruction())

	// Проверяем, что массивы содержат одинаковые данные, но не шарят общую память
	require.Len(t, clone.Blanks(), 2)
	assert.NotSame(t, &q.Blanks()[0], &clone.Blanks()[0])
}
