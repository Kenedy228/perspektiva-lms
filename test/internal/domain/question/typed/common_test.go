package typed_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/typed"
	"github.com/stretchr/testify/require"
)

func makeContent(s string) question.Content {
	opt, _ := question.NewContent(question.ContentTypeText, s)
	return opt
}

func castQuestion(t *testing.T, q question.Question) *typed.Question {
	cast, ok := q.(*typed.Question)
	require.True(t, ok)

	return cast
}

func castAnswer(t *testing.T, ans question.Answer) typed.Answer {
	cast, ok := ans.(typed.Answer)
	require.True(t, ok)

	return cast
}

func mockBlanks() []typed.Blank {
	blanks := make([]typed.Blank, 0, 5)

	for i := range 5 {
		v := fmt.Sprintf("%d", i)
		blanks = append(blanks, newBlankBuilder().withPlaceholder(v).
			withVariant(v).
			buildNoTest())
	}

	return blanks
}
