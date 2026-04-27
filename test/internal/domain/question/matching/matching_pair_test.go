package matching_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/matching"
	"github.com/stretchr/testify/assert"
)

func TestNewPair(t *testing.T) {
	t.Run("ErrEmptyPrompt", func(t *testing.T) {
		t.Run("with empty prompt", func(t *testing.T) {
			//Arrange-Assert
			newPairBuilder().withPrompt("").
				build(t, matching.ErrInvalidPrompt)
		})

		t.Run("with whitespaces prompt", func(t *testing.T) {
			//Arrange-Assert
			newPairBuilder().withPrompt(" ").
				build(t, matching.ErrInvalidPrompt)
		})
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		pair := newPairBuilder().withPrompt("my prompt").
			withContentText("content").
			build(t, nil)

		//Assert
		assert.Equal(t, pair.Prompt(), "my prompt")
		assert.Equal(t, pair.Content().ContentType(), question.ContentTypeText)
		assert.Equal(t, pair.Content().Value(), "content")
	})
}

func TestEqualPair(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		//Arrange
		first := newPairBuilder().build(t, nil)
		second := newPairBuilder().build(t, nil)

		//Assert
		assert.Equal(t, first, second)
	})

	t.Run("NotEqual", func(t *testing.T) {
		t.Run("different prompts", func(t *testing.T) {
			//Arrange
			first := newPairBuilder().withPrompt("first").
				build(t, nil)
			second := newPairBuilder().withPrompt("second").
				build(t, nil)

			//Assert
			assert.NotEqual(t, first, second)
		})

		t.Run("different option value, but same cType", func(t *testing.T) {
			//Arrange
			first := newPairBuilder().withContentText("first").
				build(t, nil)
			second := newPairBuilder().withContentText("second").
				build(t, nil)

			//Assert
			assert.NotEqual(t, first, second)
		})

		t.Run("same option value, but different cType", func(t *testing.T) {
			//Arrange
			first := newPairBuilder().withContent(question.ContentTypeImage, "same").
				build(t, nil)
			second := newPairBuilder().withContent(question.ContentTypeAudio, "same").
				build(t, nil)

			//Assert
			assert.NotEqual(t, first, second)
		})

		t.Run("different content option", func(t *testing.T) {
			//Arrange
			first := newPairBuilder().withContent(question.ContentTypeImage, "image").
				build(t, nil)
			second := newPairBuilder().withContent(question.ContentTypeAudio, "audio").
				build(t, nil)

			//Assert
			assert.NotEqual(t, first, second)
		})
	})
}
