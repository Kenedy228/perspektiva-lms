package quiz_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz"
	"gitflic.ru/lms/internal/domain/shared/limit"
	"github.com/stretchr/testify/assert"
)

type quizBuilder struct {
	title       string
	sources     []quiz.Source
	maxAttempts int
	timeLimit   limit.Limit
}

func newQuizBuilder() *quizBuilder {
	return &quizBuilder{
		title:       "",
		sources:     []quiz.Source{},
		maxAttempts: -1,
		timeLimit:   makeLimit(0),
	}
}

func (b *quizBuilder) withTitle(s string) *quizBuilder {
	b.title = s
	return b
}

func (b *quizBuilder) withSource(src quiz.Source) *quizBuilder {
	b.sources = append(b.sources, src)
	return b
}

func (b *quizBuilder) withSourceList(src []quiz.Source) *quizBuilder {
	for i := range src {
		b.sources = append(b.sources, src[i])
	}
	return b
}

func (b *quizBuilder) withMaxAttempts(val int) *quizBuilder {
	b.maxAttempts = val
	return b
}

func (b *quizBuilder) withTimeLimit(val int) *quizBuilder {
	b.timeLimit = makeLimit(val)
	return b
}

func (b *quizBuilder) build(t *testing.T, wantErr error) *quiz.Quiz {
	params := quiz.Params{
		Title:       b.title,
		Sources:     b.sources,
		MaxAttempts: b.maxAttempts,
		TimeLimit:   b.timeLimit,
	}

	q, err := quiz.New(params)
	assert.ErrorIs(t, err, wantErr)

	return q
}
