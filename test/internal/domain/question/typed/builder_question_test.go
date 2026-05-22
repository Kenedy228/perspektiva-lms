//go:build legacy
// +build legacy

package typed_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/typed"
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/question/title"
	"github.com/stretchr/testify/require"
)

type questionBuilder struct {
	titleStr string
	blanks   []blank.Blank
}

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{}
}

func (b *questionBuilder) withTitle(s string) *questionBuilder {
	b.titleStr = s
	return b
}

func (b *questionBuilder) withBlank(placeholder, variant string) *questionBuilder {
	b.blanks = append(b.blanks, makeBlank(placeholder, variant))
	return b
}

// build используется, когда мы ожидаем успешное создание агрегата
func (b *questionBuilder) build(t *testing.T) *typed.Question {
	t.Helper()
	q, err := b.buildWithError()
	require.NoError(t, err)
	return q
}

// buildWithError используется, когда мы хотим протестировать ошибки валидации
func (b *questionBuilder) buildWithError() (*typed.Question, error) {
	tTitle, err := title.New(makeContent(b.titleStr))
	if err != nil {
		return nil, err
	}
	return typed.New(tTitle, b.blanks)
}
