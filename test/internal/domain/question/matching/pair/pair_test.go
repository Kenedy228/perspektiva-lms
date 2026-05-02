package pair_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/matching/pair"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew_Success(t *testing.T) {
	//Arrange
	p, err := newPairBuilder().
		withPrompt(content.TypeText, "prompt").
		withMatch(content.TypeImage, "image.png").
		build()

	//Assert
	assert.NoError(t, err)
	assert.Equal(t, p.Prompt().Value(), "prompt")
	assert.Equal(t, p.Match().Value(), "image.png")
	assert.NotEqual(t, p.Prompt().ID(), uuid.Nil)
	assert.NotEqual(t, p.Match().ID(), uuid.Nil)
}

func TestNew_Fail(t *testing.T) {
	t.Run("некорректный prompt", func(t *testing.T) {
		t.Run("тип контента не текстовый", func(t *testing.T) {
			//Arrange
			_, err := newPairBuilder().
				withPrompt(content.TypeImage, "image.png").
				withMatch(content.TypeImage, "image.png").
				build()

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, pair.ErrInvalid)
		})

		t.Run("выход за лимит максимального количества символов", func(t *testing.T) {
			//Arrange
			_, err := newPairBuilder().
				withPrompt(content.TypeText, strings.Repeat("A", 1e5)).
				withMatch(content.TypeImage, "image.png").
				build()

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, pair.ErrInvalid)
		})
	})

	t.Run("некорректный match", func(t *testing.T) {
		t.Run("текстовый match и выход за лимит максимального количества символов", func(t *testing.T) {
			//Arrange
			_, err := newPairBuilder().
				withPrompt(content.TypeText, "prompt").
				withMatch(content.TypeText, strings.Repeat("a", 1e5)).
				build()

			//Assert
			assert.Error(t, err)
			assert.ErrorIs(t, err, pair.ErrInvalid)
		})
	})
}
