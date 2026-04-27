package question_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/stretchr/testify/assert"
)

func TestIsValidContentType(t *testing.T) {
	t.Run("for existing values return true", func(t *testing.T) {
		//Assert
		assert.True(t, question.ContentTypeText.IsValid())
		assert.True(t, question.ContentTypeAudio.IsValid())
		assert.True(t, question.ContentTypeImage.IsValid())
	})

	t.Run("for unexisting values return false", func(t *testing.T) {
		//Assert
		assert.False(t, question.ContentType("").IsValid())
		assert.False(t, question.ContentType(" ").IsValid())
	})
}

func TestTitleContentType(t *testing.T) {
	t.Run("for existing values return non-empty text", func(t *testing.T) {
		//Assert
		assert.NotEmpty(t, question.ContentTypeText.Title())
		assert.NotEmpty(t, question.ContentTypeAudio.Title())
		assert.NotEmpty(t, question.ContentTypeImage.Title())
	})

	t.Run("for unexisting values return empty text", func(t *testing.T) {
		//Assert
		assert.Empty(t, question.ContentType("").Title())
		assert.Empty(t, question.ContentType(" ").Title())
	})
}
