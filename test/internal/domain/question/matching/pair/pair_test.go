//go:build legacy
// +build legacy

package pair_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/question/content"
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
	assert.Equal(t, p.prompt().Value(), "prompt")
	assert.Equal(t, p.match().Value(), "image.png")
	assert.NotEqual(t, p.prompt().ID(), uuid.Nil)
	assert.NotEqual(t, p.match().ID(), uuid.Nil)
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
