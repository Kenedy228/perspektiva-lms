package item_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/attempt/item"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockQuestion - простая заглушка для интерфейса question.Question

func TestNew(t *testing.T) {
	t.Run("передан nil вместо вопроса возвращает ошибку", func(t *testing.T) {
		// Act
		_, err := item.New(nil)

		// Assert
		assert.ErrorIs(t, err, item.ErrInvalid)
	})

	t.Run("успешное создание", func(t *testing.T) {
		// Arrange
		qID := uuid.New()
		q := mockQuestion{id: qID}

		// Act
		itm, err := item.New(q)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, qID, itm.ID())
		assert.NotNil(t, itm.Snapshot())
	})
}
