package sequence_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnswer(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			//Arrange
			ans := castAnswer(t, newAnswerBuilder().build())

			//Assert
			assert.Equal(t, len(ans.IDs()), 0)
			assert.True(t, ans.IsEmpty())
		})

		t.Run("non empty", func(t *testing.T) {
			//Arrange
			ans := castAnswer(t, newAnswerBuilder().withRandomID().
				withRandomID().
				build())

			//Assert
			assert.Equal(t, len(ans.IDs()), 2)
			assert.False(t, ans.IsEmpty())
		})
	})
}

func TestClone(t *testing.T) {
	//Arrange
	ans := castAnswer(t, newAnswerBuilder().withRandomID().
		withRandomID().
		build())
	clone := castAnswer(t, ans.Clone())

	//Assert
	assert.Equal(t, len(ans.IDs()), len(clone.IDs()))
	assert.NotSame(t, &ans.IDs()[0], &clone.IDs()[0])
}
