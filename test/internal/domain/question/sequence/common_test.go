package sequence_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/sequence"
	"github.com/stretchr/testify/require"
)

func makeContent(s string) question.Content {
	opt, _ := question.NewContent(question.ContentTypeText, s)
	return opt
}

func castQuestion(t *testing.T, q question.Question) *sequence.Question {
	cast, ok := q.(*sequence.Question)
	require.True(t, ok)

	return cast
}

func castAnswer(t *testing.T, ans question.Answer) sequence.Answer {
	cast, ok := ans.(sequence.Answer)
	require.True(t, ok)

	return cast
}
