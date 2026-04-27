package short_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnswer(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("not empty", func(t *testing.T) {
			//Arrange
			answer := castAnswer(t, newAnswerBuilder().withInput("item").
				build())

			//Assert
			assert.Equal(t, answer.Input(), "item")
			assert.False(t, answer.IsEmpty())
		})
	})
}

func TestClone(t *testing.T) {
	//Arrange
	ans := castAnswer(t, newAnswerBuilder().build())
	clone := castAnswer(t, ans.Clone())

	//Assert
	assert.Equal(t, ans.Input(), clone.Input())
	assert.NotSame(t, &ans, &clone)
}
