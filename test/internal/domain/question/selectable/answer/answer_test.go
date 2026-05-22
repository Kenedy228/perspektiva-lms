//go:build legacy
// +build legacy

package answer_test

import (
	"testing"

	answer2 "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	ans := answer2.New(makeOptions(10))

	//Assert
	assert.Equal(t, len(ans.Options()), 10)
	assert.Equal(t, len(ans.OptionsAsMap()), 10)
}

func TestIsEmpty(t *testing.T) {
	t.Run("при пустом слайсе возвращает true", func(t *testing.T) {
		//Arrange
		ans := answer2.New(makeOptions(0))

		//Assert
		assert.True(t, ans.IsEmpty())
	})

	t.Run("при непустом слайсе возвращает false", func(t *testing.T) {
		//Arrange
		ans := answer2.New(makeOptions(10))

		//Assert
		assert.False(t, ans.IsEmpty())
	})
}

func TestOptions(t *testing.T) {
	//Arrange
	ans := answer2.New(makeOptions(10))
	opts := ans.Options()

	//Act
	opts[0] = answer.AnswerOption{
		OptionID: uuid.New(),
	}

	//Assert
	assert.NotEqual(t, opts[0], ans.Options()[0])
}

func TestClone(t *testing.T) {
	//Arrange
	ans := answer2.New(makeOptions(10))
	clone, ok := ans.Clone().(answer2.Answer)
	require.True(t, ok)

	//Assert
	assert.Equal(t, len(ans.Options()), len(clone.Options()))
	assert.Equal(t, ans.Options(), clone.Options())
}
