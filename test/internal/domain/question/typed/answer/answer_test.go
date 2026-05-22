//go:build legacy
// +build legacy

package answer_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/typed/answer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	ans := answer.New(makeAnswerBlanks(10))

	//Assert
	assert.Equal(t, 10, len(ans.Blanks()))
}

func TestIsEmpty(t *testing.T) {
	tc := []struct {
		name        string
		blanksCount int
		isEmpty     bool
	}{
		{
			name:        "если пустой срез blanks - возвращает true",
			blanksCount: 0,
			isEmpty:     true,
		},
		{
			name:        "если непустой срез blanks - возврщает false",
			blanksCount: 10,
			isEmpty:     false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			ans := answer.New(makeAnswerBlanks(tt.blanksCount))

			//Assert
			assert.Equal(t, tt.isEmpty, ans.IsEmpty())
		})
	}
}

func TestClone(t *testing.T) {
	//Arrange
	ans := answer.New(makeAnswerBlanks(10))
	clone, ok := ans.Clone().(answer.Answer)
	require.True(t, ok)

	//Assert
	assert.Equal(t, ans.Blanks(), clone.Blanks())
}
