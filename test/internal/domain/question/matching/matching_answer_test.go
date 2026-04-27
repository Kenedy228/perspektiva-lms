package matching_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnswer(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			//Arrange
			ans := newAnswerBuilder().build(t)

			//Assert
			assert.True(t, ans.IsEmpty())
		})

		t.Run("non-empty", func(t *testing.T) {
			//Arrange
			ans := newAnswerBuilder().withDefaultPair().
				build(t)

			//Assert
			assert.False(t, ans.IsEmpty())
		})
	})
}

func TestAnswerClone(t *testing.T) {
	t.Run("should clone answer", func(t *testing.T) {
		//Arrange
		ans := castAnswer(t, newAnswerBuilder().withDefaultPair().
			build(t))
		clone := castAnswer(t, ans.Clone())

		//Assert
		assert.NotSame(t, &ans, &clone)
		assert.NotSame(t, &clone.Pairs()[0], &clone.Pairs()[0])
		assert.Equal(t, &clone.Pairs()[0], &clone.Pairs()[0])
	})
}
