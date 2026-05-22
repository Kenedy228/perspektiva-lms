//go:build legacy
// +build legacy

package quiz_test

import (
	"testing"

	quiz2 "gitflic.ru/lms/backend/internal/domain/quiz"
	limit2 "gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
	"github.com/stretchr/testify/assert"
)

type quizBuilder struct {
	title       title.Title
	sources     []source.Source
	maxAttempts limit2.Attempts
	timeLimit   limit2.Time
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

func (b *quizBuilder) build(t *testing.T, wantErr error) *quiz2.Quiz {
	t.Helper()

	params := quiz2.Params{
		Title:       b.title,
		Sources:     b.sources,
		MaxAttempts: b.maxAttempts,
		TimeLimit:   b.timeLimit,
	}

	q, err := quiz2.New(params)
	if wantErr != nil {
		assert.ErrorIs(t, err, wantErr)
		return nil
	}

	assert.NoError(t, err)
	return q
}
