package typed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAnswer(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			//Arrange
			answer := castAnswer(t, newAnswerBuilder().withBlank("{{placeholder}}", "value").
				build())

			//Assert
			assert.Contains(t, answer.Inputs(), "{{placeholder}}")
			assert.False(t, answer.IsEmpty())
		})

		t.Run("not empty", func(t *testing.T) {
			//Arrange
			answer := castAnswer(t, newAnswerBuilder().build())

			//Assert
			assert.True(t, answer.IsEmpty())
		})
	})
}

func TestCloneAnswer(t *testing.T) {
	//Arrange
	a := castAnswer(t, newAnswerBuilder().withBlank("{{placeholder}}", "value").build())
	clone := castAnswer(t, a.Clone())

	//Assert
	assert.Contains(t, a.Inputs(), "{{placeholder}}")
	assert.Contains(t, clone.Inputs(), "{{placeholder}}")
}
