package matching_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/matching"
)

func newAnswerBuilder() *answerBuilder {
	return &answerBuilder{
		pairs: []matching.Pair{},
	}
}

type answerBuilder struct {
	pairs []matching.Pair
}

func (b *answerBuilder) withPair(prompt string, val string, cType question.ContentType) *answerBuilder {
	pair := newPairBuilder().withPrompt(prompt).withContent(cType, val).buildNoTest()
	b.pairs = append(b.pairs, pair)
	return b
}

func (b *answerBuilder) withTextPair(prompt string, val string) *answerBuilder {
	b.withPair(prompt, val, question.ContentTypeText)
	return b
}

func (b *answerBuilder) withDefaultPair() *answerBuilder {
	pair := newPairBuilder().buildNoTest()
	b.pairs = append(b.pairs, pair)
	return b
}

func (b *answerBuilder) build(t *testing.T) question.Answer {
	t.Helper()

	params := matching.AnswerParams{
		Pairs: b.pairs,
	}

	ans := matching.NewAnswer(params)
	return ans
}
