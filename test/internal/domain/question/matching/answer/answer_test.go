package answer_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/matching/answer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	pairs := []answer.AnswerPair{
		makeAnswerPair(uuid.New(), uuid.New()),
	}
	ans := answer.New(pairs)

	//Assert
	assert.Equal(t, 1, len(ans.Pairs()))
	assert.Equal(t, 1, len(ans.PairsAsMap()))
}

func TestIsEmpty(t *testing.T) {
	t.Run("при пустых ответах возвращает true", func(t *testing.T) {
		//Arrange
		ans := answer.New([]answer.AnswerPair{})

		//Assert
		assert.True(t, ans.IsEmpty())
	})

	t.Run("при непустых ответах возвращает false", func(t *testing.T) {
		//Arrange
		pairs := []answer.AnswerPair{
			makeAnswerPair(uuid.New(), uuid.New()),
		}
		ans := answer.New(pairs)

		//Assert
		assert.False(t, ans.IsEmpty())
	})
}

func TestAnswerClone(t *testing.T) {
	//Arrange
	pairs := []answer.AnswerPair{
		makeAnswerPair(uuid.New(), uuid.New()),
	}
	ans := answer.New(pairs)
	clone, ok := ans.Clone().(answer.Answer)
	require.True(t, ok)

	//Assert
	assert.Equal(t, len(clone.Pairs()), len(ans.Pairs()))
}
