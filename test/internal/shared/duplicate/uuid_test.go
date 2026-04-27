package duplicate_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/shared/duplicate"
	"github.com/stretchr/testify/assert"
)

func TestFindUUID(t *testing.T) {
	t.Run("should return false if there is no duplicates", func(t *testing.T) {
		//Arrange
		ids := makeIDs(1000)

		//Act
		res := duplicate.FindUUID(ids)

		//Assert
		assert.False(t, res)
	})

	t.Run("should return true if there is duplciates", func(t *testing.T) {
		//Arrange
		ids := makeIDs(1000)
		ids = append(ids, ids[0])

		//Act
		res := duplicate.FindUUID(ids)

		//Assert
		assert.True(t, res)
	})
}

func BenchmarkFindUUID(b *testing.B) {
	sizes := []int{1e3, 1e4, 1e5}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size %d", size), func(b *testing.B) {
			ids := makeIDs(size)

			for b.Loop() {
				duplicate.FindUUID(ids)
			}
		})
	}
}
