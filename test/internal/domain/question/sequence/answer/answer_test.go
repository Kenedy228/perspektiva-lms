//go:build legacy
// +build legacy

package answer_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	ans := answer.New(makeAnswerOptions(10))

	//Assert
	assert.Equal(t, 10, len(ans.Options()))
	assert.False(t, ans.IsEmpty())
}

func TestIsEmpty(t *testing.T) {
	tc := []struct {
		name         string
		optionsCount int
		isEmpty      bool
	}{
		{
			name:         "если нет опций - возвращает true",
			optionsCount: 0,
			isEmpty:      true,
		},
		{
			name:         "если есть опции - возвращает false",
			optionsCount: 10,
			isEmpty:      false,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			ans := answer.New(makeAnswerOptions(10))

			//Assert
			assert.Equal(t, tt.isEmpty, ans.IsEmpty())
		})
	}
}

func TestClone(t *testing.T) {
	//Arrange
	ans := answer.New(makeAnswerOptions(10))
	clone, ok := ans.Clone().(answer.Answer)
	require.True(t, ok)

	//Assert
	assert.Equal(t, ans.Options(), clone.Options())
	assert.Equal(t, ans.IsEmpty(), clone.IsEmpty())
}
