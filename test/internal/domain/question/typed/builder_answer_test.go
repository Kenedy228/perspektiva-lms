package typed_test

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/typed"
)

func newAnswerBuilder() *answerBuilder {
	return &answerBuilder{
		items: map[string]string{},
	}
}

type answerBuilder struct {
	items map[string]string
}

func (b *answerBuilder) withBlank(p, v string) *answerBuilder {
	b.items[p] = v
	return b
}

func (b *answerBuilder) build() question.Answer {
	params := typed.AnswerParams{
		Inputs: b.items,
	}

	return typed.NewAnswer(params)
}
