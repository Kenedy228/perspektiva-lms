package question_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	t.Run("return true for existing consts", func(t *testing.T) {
		//Assert
		assert.True(t, question.TypeMatching.IsValid())
		assert.True(t, question.TypeSelectable.IsValid())
		assert.True(t, question.TypeSequence.IsValid())
		assert.True(t, question.TypeTyped.IsValid())
		assert.True(t, question.TypeShort.IsValid())
	})

	t.Run("return false for non-existing consts", func(t *testing.T) {
		//Assert
		assert.False(t, question.Type("").IsValid())
		assert.False(t, question.Type("foo").IsValid())
	})
}

func TestTitle(t *testing.T) {
	t.Run("return non-empty string for existing consts", func(t *testing.T) {
		//Assert
		assert.NotEmpty(t, question.TypeMatching.Title())
		assert.NotEmpty(t, question.TypeSelectable.Title())
		assert.NotEmpty(t, question.TypeSequence.Title())
		assert.NotEmpty(t, question.TypeTyped.Title())
		assert.NotEmpty(t, question.TypeShort.Title())
	})

	t.Run("return empty string for non-existing consts", func(t *testing.T) {
		//Assert
		assert.Empty(t, question.Type("").Title())
		assert.Empty(t, question.Type("foo").Title())
	})
}

func TestDefaultInstruction(t *testing.T) {
	t.Run("return non-empty string for existing consts", func(t *testing.T) {
		//Assert
		assert.NotEmpty(t, question.TypeMatching.DefaultInstruction())
		assert.NotEmpty(t, question.TypeSelectable.DefaultInstruction())
		assert.NotEmpty(t, question.TypeSequence.DefaultInstruction())
		assert.NotEmpty(t, question.TypeTyped.DefaultInstruction())
		assert.NotEmpty(t, question.TypeShort.DefaultInstruction())
	})

	t.Run("return empty string for non-existing consts", func(t *testing.T) {
		//Assert
		assert.Empty(t, question.Type("").DefaultInstruction())
		assert.Empty(t, question.Type("foo").DefaultInstruction())
	})
}
