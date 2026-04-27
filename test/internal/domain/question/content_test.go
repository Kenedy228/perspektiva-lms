package question_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/stretchr/testify/assert"
)

func TestNewContent(t *testing.T) {
	t.Run("should return error if invalid content type", func(t *testing.T) {
		//Arrange-Assert
		newContentBuilder().withValue("value").build(t, question.ErrInvalidContent)
	})

	t.Run("should return error for values", func(t *testing.T) {
		t.Run("should return error if content type text and value has no non-whitespace symbols", func(t *testing.T) {
			//Arrange-Assert
			newContentBuilder().withTextCType().withValue("").build(t, question.ErrInvalidContent)
			newContentBuilder().withTextCType().withValue(" ").build(t, question.ErrInvalidContent)
		})

		t.Run("sgould return error if content type image or audio and value is not valid s3 key", func(t *testing.T) {
			//Arrange-Assert
			newContentBuilder().withImageCType().withValue(" key 123").build(t, question.ErrInvalidContent)
			newContentBuilder().withAudioCType().withValue(" key 123").build(t, question.ErrInvalidContent)
		})
	})

	t.Run("should return no err", func(t *testing.T) {
		//Arrange
		c := newContentBuilder().withTextCType().withValue("value").build(t, nil)

		//Assert
		assert.Equal(t, c.ContentType(), question.ContentTypeText)
		assert.Equal(t, c.Value(), "value")
	})
}

func TestIsText(t *testing.T) {
	t.Run("for non text content types return false", func(t *testing.T) {
		//Arrange
		image := newContentBuilder().withImageCType().withValue("image").build(t, nil)
		audio := newContentBuilder().withAudioCType().withValue("audio").build(t, nil)

		//Assert
		assert.False(t, image.IsText())
		assert.False(t, audio.IsText())
	})

	t.Run("for text content type return true", func(t *testing.T) {
		//Arrange
		text := newContentBuilder().withTextCType().withValue("text").build(t, nil)

		//Assert
		assert.True(t, text.IsText())
	})
}

func TestEqualContent(t *testing.T) {
	t.Run("not equal", func(t *testing.T) {
		t.Run("diff cType, equal values", func(t *testing.T) {
			//Arrange
			text := newContentBuilder().withTextCType().withValue("value").build(t, nil)
			image := newContentBuilder().withImageCType().withValue("value").build(t, nil)

			//Assert
			assert.False(t, text.Equal(image))
		})

		t.Run("equal cType, diff values", func(t *testing.T) {
			//Arrange
			firstText := newContentBuilder().withTextCType().withValue("value").build(t, nil)
			secondText := newContentBuilder().withTextCType().withValue("value2").build(t, nil)

			//Assert
			assert.False(t, firstText.Equal(secondText))
		})
	})

	t.Run("equal", func(t *testing.T) {
		//Arrange
		firstText := newContentBuilder().withTextCType().withValue("value").build(t, nil)
		secondText := newContentBuilder().withTextCType().withValue("value").build(t, nil)

		//Assert
		assert.True(t, firstText.Equal(secondText))
	})
}
