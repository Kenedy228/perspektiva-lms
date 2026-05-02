package limit_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/limit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("должен вернуть ошибку, если длительность отрицательная", func(t *testing.T) {
		//Arrange
		_, err := limit.NewTime(-1)

		//Assert
		assert.Error(t, err)
	})

	t.Run("должен вернуть ошибку, если длительность больше разрешенного лимита", func(t *testing.T) {
		//Arrange
		_, err := limit.NewTime(1e5)

		//Assert
		assert.Error(t, err)
	})

	t.Run("должен создать объект, если значение в пределах границ (см. timelimit)", func(t *testing.T) {
		t.Run("бесконечный лимит", func(t *testing.T) {
			//Arrange
			limit, err := limit.NewTime(0)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, limit.Seconds(), 0)
			assert.True(t, limit.IsInfinite())
		})

		t.Run("конечный лимит", func(t *testing.T) {
			//Arrange
			limit, err := limit.NewTime(5000)

			//Assert
			assert.NoError(t, err)
			assert.Equal(t, limit.Seconds(), 5000)
			assert.False(t, limit.IsInfinite())
		})
	})
}

func TestTryDuration(t *testing.T) {
	t.Run("должен вернуть объект time.Duration и значение true, если лимит не бесконечный", func(t *testing.T) {
		//Arrange
		limit, err := limit.NewTime(5000)
		require.NoError(t, err)

		//Act
		duration, ok := limit.TryDuration()

		//Assert
		assert.Equal(t, limit.Seconds(), 5000)
		assert.False(t, limit.IsInfinite())
		assert.Equal(t, duration.Seconds(), 5000.0)
		assert.True(t, ok)
	})

	t.Run("должен вернуть объект time.Duration и значение false, если лимит бесконечный", func(t *testing.T) {
		//Arrange
		limit, err := limit.NewTime(0)
		require.NoError(t, err)

		//Act
		duration, ok := limit.TryDuration()

		//Assert
		assert.Equal(t, limit.Seconds(), 0)
		assert.True(t, limit.IsInfinite())
		assert.Equal(t, duration.Seconds(), 0.0)
		assert.False(t, ok)
	})
}
