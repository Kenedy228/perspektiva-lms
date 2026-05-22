//go:build legacy
// +build legacy

package typed_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/typed"
	answer2 "gitflic.ru/lms/backend/internal/domain/question/typed/answer"
	"gitflic.ru/lms/backend/internal/domain/question/typed/blank"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"gitflic.ru/lms/backend/internal/domain/question/title"
	"github.com/stretchr/testify/require"
)

func makeAnswer(blanks map[string]string) answer2.Answer {
	var ansBlanks []answer2.AnswerBlank
	for placeholder, variant := range blanks {
		ansBlanks = append(ansBlanks, answer2.AnswerBlank{
			Placeholder: placeholder,
			Variant:     variant,
		})
	}
	return answer2.New(ansBlanks)
}

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

func (b *questionBuilder) withBlank(placeholder string, variantVals ...string) *questionBuilder {
	var contents []content.Content
	for _, val := range variantVals {
		contents = append(contents, makeContent(val))
	}

	newBlank, _ := blank.New(placeholder, contents)
	b.blanks = append(b.blanks, newBlank)
	return b
}

func (b *questionBuilder) build(t *testing.T) *typed.Question {
	t.Helper()

	tTitle, err := title.New(makeContent(b.titleStr))
	require.NoError(t, err)

	q, err := typed.New(tTitle, b.blanks)
	require.NoError(t, err)

	return q
}

func makeContent(value string) content.Content {
	c, _ := content.New(content.TypeText, value)
	return c
}
