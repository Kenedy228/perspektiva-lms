package answer_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/short/answer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	//Arrange
	ans := answer.New(answer.AnswerVariant{
		Input: "input",
	})

	//Assert
	assert.Equal(t, "input", ans.VariantAsString())
}

func TestIsEmpty(t *testing.T) {
	tc := []struct {
		name    string
		val     string
		isEmpty bool
	}{
		{
			name:    "непустая строка - непустое значение",
			val:     "текст ответа",
			isEmpty: false,
		},
		{
			name:    "пустая строка - пустое значение",
			val:     "",
			isEmpty: true,
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			//Arrange
			ans := answer.New(answer.AnswerVariant{
				Input: "input",
			})

			//Assert
			assert.Equal(t, tt.isEmpty, ans.IsEmpty())
		})
	}
}

func TestClone(t *testing.T) {
	//Arrange
	ans := answer.New(answer.AnswerVariant{
		Input: "input",
	})
	clone, ok := ans.Clone().(answer.Answer)
	require.True(t, ok)

	//Assert
	assert.Equal(t, ans.VariantAsString(), clone.VariantAsString())
}
