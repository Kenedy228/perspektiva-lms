//go:build legacy
// +build legacy

package uid_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("должен вернуть идентификатор без ошибки", func(t *testing.T) {
		//Arrange
		id, err := uid.New()

		//Assert
		assert.NoError(t, err)
		assert.NotEqual(t, id, uuid.Nil)
	})

	t.Run("строковое представление возвращаемого идентификатора должно парситься в валидный uuid", func(t *testing.T) {
		//Arrange
		id, err := uid.New()
		require.NoError(t, err)

		//Act
		parsed, err := uuid.Parse(id.String())

		//Assert
		assert.NoError(t, err)
		assert.NotEqual(t, parsed, uuid.Nil)
	})

	t.Run("разные вызовы конструктора возвращают разные идентификаторы", func(t *testing.T) {
		//Arrange
		first, err := uid.New()
		require.NoError(t, err)

		second, err := uid.New()
		require.NoError(t, err)

		//Assert
		assert.NotEqual(t, first, second)
		assert.NotEqual(t, first, uuid.Nil)
		assert.NotEqual(t, second, uuid.Nil)
	})
}
