//go:build legacy
// +build legacy

package sequence_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/sequence"
	answer2 "gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	"gitflic.ru/lms/backend/internal/domain/question/title"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func makeAnswer(optionIDs []uuid.UUID) answer2.Answer {
	var ansOptions []answer.AnswerOption
	for _, id := range optionIDs {
		ansOptions = append(ansOptions, answer.AnswerOption{
			OptionID: id,
		})
	}
	return answer2.New(ansOptions)
}

type questionBuilder struct {
	titleStr string
	options  []option.Option
}

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{}
}

func (b *questionBuilder) withTitle(s string) *questionBuilder {
	b.titleStr = s
	return b
}

func (b *questionBuilder) withOption(variantVal string) *questionBuilder {
	c, _ := content.New(content.TypeText, variantVal)
	opt, _ := option.New(c)
	b.options = append(b.options, opt)
	return b
}

func (b *questionBuilder) build(t *testing.T) *sequence.Question {
	t.Helper()

	tTitle, err := title.New(makeContent(b.titleStr))
	require.NoError(t, err)

	q, err := sequence.New(tTitle, b.options)
	require.NoError(t, err)

	return q
}

func makeContent(value string) content.Content {
	c, _ := content.New(content.TypeText, value)
	return c
}
