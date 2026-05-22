//go:build legacy
// +build legacy

package criteria_test

import (
	"testing"

	criteria2 "gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/stretchr/testify/assert"
)

func TestNewRandom(t *testing.T) {
	t.Run("при отрицательном значении возвращает error", func(t *testing.T) {
		//Arrange-Assert
		newRandomBuilder().build(t, criteria2.ErrInvalid)
		newRandomBuilder().withCount(-1).
			build(t, criteria2.ErrInvalid)
	})

	t.Run("при выходе за лимиты возвращает error", func(t *testing.T) {
		//Arrange-Assert
		newRandomBuilder().withCount(1e9).build(t, criteria2.ErrInvalid)
	})

	t.Run("корректное количество", func(t *testing.T) {
		//Arrange
		c := castRandom(t, newRandomBuilder().withCount(500).build(t, nil))

		//Assert
		assert.Equal(t, c.Type(), criteria2.TypeRandom)
		assert.Equal(t, c.QuestionCount(), 500)
	})
}
