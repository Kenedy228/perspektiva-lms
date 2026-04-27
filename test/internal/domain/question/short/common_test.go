package short_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/short"
	"github.com/stretchr/testify/require"
)

func makeContent(s string) question.Content {
	opt, _ := question.NewContent(question.ContentTypeText, s)
	return opt
}

func castQuestion(t *testing.T, q question.Question) *short.Question {
	cast, ok := q.(*short.Question)
	require.True(t, ok)

	return cast
}

func castAnswer(t *testing.T, ans question.Answer) short.Answer {
	cast, ok := ans.(short.Answer)
	require.True(t, ok)

	return cast
}

func mockVariants() []question.Content {
	items := make([]question.Content, 0, 5)

	for i := range 5 {
		items = append(items, makeContent(fmt.Sprintf("%d", i)))
	}

	return items
}
