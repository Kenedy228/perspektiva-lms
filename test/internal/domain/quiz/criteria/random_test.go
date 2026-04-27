package criteria_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/stretchr/testify/assert"
)

func TestNewRandom(t *testing.T) {
	t.Run("for count <= 0 returns error", func(t *testing.T) {
		//Arrange-Assert
		newRandomBuilder().build(t, criteria.ErrInvalidQuestionCount)
		newRandomBuilder().withCount(-1).
			build(t, criteria.ErrInvalidQuestionCount)
	})

	t.Run("for positive count greater than maxQuestions returns error", func(t *testing.T) {
		//Arrange-Assert
		newRandomBuilder().withCount(1e9).build(t, criteria.ErrInvalidQuestionCount)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		c := castRandom(t, newRandomBuilder().withCount(500).build(t, nil))
		
		//Assert
		assert.Equal(t, c.Type(), criteria.TypeRandom)
		assert.Equal(t, c.QuestionCount(), 500)
	})
}
