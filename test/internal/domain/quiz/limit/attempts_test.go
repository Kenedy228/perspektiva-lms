//go:build legacy
// +build legacy

package limit_test

import (
	"testing"

	limit2 "gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAttempts(t *testing.T) {
	t.Run("должен вернуть ошибку, если количество попыток отрицательное", func(t *testing.T) {
		// Arrange
		_, err := limit2.NewAttempts(-1)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, limit2.ErrInvalid)
		assert.Contains(t, err.Error(), "не может быть отрицательным")
	})

	t.Run("должен вернуть ошибку, если количество попыток больше разрешенного лимита", func(t *testing.T) {
		// Arrange
		_, err := limit2.NewAttempts(1e5)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, limit2.ErrInvalid)
		assert.Contains(t, err.Error(), "не должно превышать")
	})

	t.Run("должен создать объект, если значение в пределах границ", func(t *testing.T) {
		t.Run("бесконечное количество попыток (ноль)", func(t *testing.T) {
			// Arrange
			att, err := limit2.NewAttempts(0)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, 0, att.Count())
			assert.True(t, att.IsInfinite())
		})

		t.Run("конечный лимит (обычный)", func(t *testing.T) {
			// Arrange
			att, err := limit2.NewAttempts(3)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, 3, att.Count())
			assert.False(t, att.IsInfinite())
		})
	})
}
