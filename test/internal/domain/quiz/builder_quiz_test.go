package quiz_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz"
	"gitflic.ru/lms/internal/domain/quiz/limit"
	"gitflic.ru/lms/internal/domain/quiz/source"
	"gitflic.ru/lms/internal/domain/quiz/title"
	"github.com/stretchr/testify/assert"
)

type quizBuilder struct {
	title       title.Title
	sources     []source.Source
	maxAttempts limit.Attempts
	timeLimit   limit.Time
}

func newQuizBuilder() *quizBuilder {
	return &quizBuilder{
		title:       makeTitle("Default Title"),
		sources:     []source.Source{},
		maxAttempts: makeAttempts(0),
		timeLimit:   makeTime(0),
	}
}

func (b *quizBuilder) withTitle(s string) *quizBuilder {
	b.title = makeTitle(s)
	return b
}

func (b *quizBuilder) withSource(src source.Source) *quizBuilder {
	b.sources = append(b.sources, src)
	return b
}

func (b *quizBuilder) withSourceList(src []source.Source) *quizBuilder {
	b.sources = append(b.sources, src...)
	return b
}

func (b *quizBuilder) withMaxAttempts(val int) *quizBuilder {
	b.maxAttempts = makeAttempts(val)
	return b
}

func (b *quizBuilder) withTimeLimit(val int) *quizBuilder {
	b.timeLimit = makeTime(val)
	return b
}

func (b *quizBuilder) build(t *testing.T, wantErr error) *quiz.Quiz {
	t.Helper()

	params := quiz.Params{
		Title:       b.title,
		Sources:     b.sources,
		MaxAttempts: b.maxAttempts,
		TimeLimit:   b.timeLimit,
	}

	q, err := quiz.New(params)
	if wantErr != nil {
		assert.ErrorIs(t, err, wantErr)
		return nil
	}

	assert.NoError(t, err)
	return q
}
