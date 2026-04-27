package short_test

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/short"
)

func newAnswerBuilder() *answerBuilder {
	return &answerBuilder{
		input: "text",
	}
}

type answerBuilder struct {
	input string
}

func (b *answerBuilder) withInput(s string) *answerBuilder {
	b.input = s
	return b
}

func (b *answerBuilder) build() question.Answer {
	params := short.AnswerParams{
		Input: b.input,
	}

	return short.NewAnswer(params)
}
