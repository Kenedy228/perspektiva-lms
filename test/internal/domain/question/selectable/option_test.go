package selectable_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewOption(t *testing.T) {
	t.Run("valid params", func(t *testing.T) {
		//Arrange
		option := newOptionBuilder().withContentAsText("content").
			withCorrect(true).
			build(t, nil)

		//Assert
		assert.NotEqual(t, option.ID(), uuid.Nil)
		assert.Equal(t, option.Content(), makeContent("content"))
		assert.True(t, option.IsCorrect())
	})
}

func TestEqualOption(t *testing.T) {
	t.Run("equal, same correct", func(t *testing.T) {
		//Arrange
		first := newOptionBuilder().withContentAsText("text").
			withCorrect(true).
			build(t, nil)
		second := newOptionBuilder().withContentAsText("text").
			withCorrect(true).
			build(t, nil)

		//Assert
		assert.True(t, first.Equal(second))
	})

	t.Run("equal, diff correct", func(t *testing.T) {
		//Arrange
		first := newOptionBuilder().withContentAsText("text").
			withCorrect(true).
			build(t, nil)
		second := newOptionBuilder().withContentAsText("text").
			withCorrect(false).
			build(t, nil)

		//Assert
		assert.True(t, first.Equal(second))
	})

	t.Run("not equal, diff content value", func(t *testing.T) {
		//Arrange
		first := newOptionBuilder().withContentAsText("first").
			withCorrect(true).
			build(t, nil)
		second := newOptionBuilder().withContentAsText("second").
			withCorrect(false).
			build(t, nil)

		//Assert
		assert.False(t, first.Equal(second))
	})

	t.Run("not equal, diff content type", func(t *testing.T) {
		//Arrange
		first := newOptionBuilder().withContent(question.ContentTypeImage, "image").
			withCorrect(true).
			build(t, nil)
		second := newOptionBuilder().withContentAsText("text").
			withCorrect(false).
			build(t, nil)

		//Assert
		assert.False(t, first.Equal(second))
	})
}
