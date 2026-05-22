//go:build legacy
// +build legacy

package answer_test

import (
	"testing"
	"time"

	answer2 "gitflic.ru/lms/backend/internal/domain/attempt/answer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	validQID := uuid.New()
	validAns := mockAnswer{}
	now := time.Now()

	t.Run("пустой ID вопроса возвращает ошибку", func(t *testing.T) {
		// Act
		entry, err := answer2.New(uuid.Nil, validAns, now)

		// Assert
		assert.ErrorIs(t, err, answer2.ErrInvalid)
		assert.Equal(t, uuid.Nil, entry.QuestionID())
	})

	t.Run("пустой ответ (nil) возвращает ошибку", func(t *testing.T) {
		// Act
		_, err := answer2.New(validQID, nil, now)

		// Assert
		assert.ErrorIs(t, err, answer2.ErrInvalid)
	})

	t.Run("успешное создание", func(t *testing.T) {
		// Act
		entry, err := answer2.New(validQID, validAns, now)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, validQID, entry.QuestionID())
		assert.NotNil(t, entry.Answer())
		assert.Equal(t, now, entry.AnsweredAt())
	})
}
