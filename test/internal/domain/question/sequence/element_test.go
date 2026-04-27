package sequence_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewElement(t *testing.T) {
	t.Run("valid params", func(t *testing.T) {
		//Arrange
		element := newElementBuilder().withContentAsText("content").
			build(t, nil)

		//Assert
		assert.NotEqual(t, element.ID(), uuid.Nil)
		assert.Equal(t, element.Content(), makeContent("content"))
	})
}

func TestEqualElement(t *testing.T) {
	t.Run("equal, same correct", func(t *testing.T) {
		//Arrange
		first := newElementBuilder().withContentAsText("text").
			build(t, nil)
		second := newElementBuilder().withContentAsText("text").
			build(t, nil)

		//Assert
		assert.True(t, first.Equal(second))
	})

	t.Run("not equal, diff content value", func(t *testing.T) {
		//Arrange
		first := newElementBuilder().withContentAsText("first").
			build(t, nil)
		second := newElementBuilder().withContentAsText("second").
			build(t, nil)

		//Assert
		assert.False(t, first.Equal(second))
	})

	t.Run("not equal, diff content type", func(t *testing.T) {
		//Arrange
		first := newElementBuilder().withContent(question.ContentTypeImage, "image").
			build(t, nil)
		second := newElementBuilder().withContentAsText("text").
			build(t, nil)

		//Assert
		assert.False(t, first.Equal(second))
	})
}
