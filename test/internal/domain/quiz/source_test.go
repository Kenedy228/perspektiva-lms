package quiz_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	t.Run("if bank id is nil return err", func(t *testing.T) {
		//Arrange-Assert
		c := new(mockCriteria)
		newSourceBuilder().withCriteria(c).
			build(t, quiz.ErrInvalidSource)
	})

	t.Run("if criteria is nil return err", func(t *testing.T) {
		//Arrange-Assert
		newSourceBuilder().withBank(uuid.New()).
			build(t, quiz.ErrInvalidSource)
	})

	t.Run("valid", func(t *testing.T) {
		//Arrange
		c := new(mockCriteria)
		s := newSourceBuilder().withCriteria(c).
			withBank(uuid.New()).
			build(t, nil)

		//Assert
		assert.NotEqual(t, s.BankID(), uuid.Nil)
		assert.NotNil(t, s.Criteria)
	})
}
