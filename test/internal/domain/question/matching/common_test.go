package matching_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/matching"
	"github.com/stretchr/testify/require"
)

func defaultPairPrompt() string {
	return "prompt"
}

func defaultPairOption() question.Content {
	return makeContent("cont")
}

func makeContent(s string) question.Content {
	opt, _ := question.NewContent(question.ContentTypeText, s)
	return opt
}

func castQuestion(t *testing.T, q question.Question) *matching.Question {
	cast, ok := q.(*matching.Question)
	require.True(t, ok)
	return cast
}

func castAnswer(t *testing.T, answer question.Answer) matching.Answer {
	cast, ok := answer.(matching.Answer)
	require.True(t, ok)
	return cast
}
