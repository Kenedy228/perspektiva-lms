package duplicates_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/shared/duplicates"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHasUUID(t *testing.T) {
	t.Run("возвращает false, если массив пустой", func(t *testing.T) {
		//Act
		has := duplicates.HasUUID([]uuid.UUID{})

		//Assert
		assert.False(t, has)
	})

	t.Run("возвращает false, если в массиве не содержатся дубликаты", func(t *testing.T) {
		//Arrange
		ids := makeIDs(10000)

		//Act
		has := duplicates.HasUUID(ids)

		//Assert
		assert.False(t, has)
	})

	t.Run("возвращает true, если в массиве содержится хотя бы один дубликат", func(t *testing.T) {
		//Arrange
		ids := makeIDs(10000)
		ids = append(ids, ids[0])

		//Act
		has := duplicates.HasUUID(ids)

		//Assert
		assert.True(t, has)
	})
}

func BenchmarkHasUUID(b *testing.B) {
	sizes := []int{1e3, 1e4, 1e5}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("размер %d", size), func(b *testing.B) {
			ids := makeIDs(size)

			for b.Loop() {
				duplicates.HasUUID(ids)
			}
		})
	}
}
