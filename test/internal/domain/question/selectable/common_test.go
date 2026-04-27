package selectable_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/selectable"
	"github.com/stretchr/testify/require"
)

func makeContent(s string) question.Content {
	content, _ := question.NewContent(question.ContentTypeText, s)
	return content
}

func castQuestion(t *testing.T, q question.Question) *selectable.Question {
	cast, ok := q.(*selectable.Question)
	require.True(t, ok)

	return cast
}

func castAnswer(t *testing.T, ans question.Answer) selectable.Answer {
	cast, ok := ans.(selectable.Answer)
	require.True(t, ok)

	return cast
}
