package typed_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/typed"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{
		text:    "text",
		imageID: uuid.Nil,
		blanks:  []typed.Blank{},
	}
}

type questionBuilder struct {
	text    string
	imageID uuid.UUID
	blanks  []typed.Blank
}

func (b *questionBuilder) withText(s string) *questionBuilder {
	b.text = s
	return b
}

func (b *questionBuilder) withImage(id uuid.UUID) *questionBuilder {
	b.imageID = id
	return b
}

func (b *questionBuilder) withBlank(p string, v ...string) *questionBuilder {
	bb := newBlankBuilder().withPlaceholder(p)

	for i := range v {
		bb.withVariant(v[i])
	}
	blank := bb.buildNoTest()

	b.blanks = append(b.blanks, blank)
	return b
}

func (b *questionBuilder) build(t *testing.T, wantErr error) question.Question {
	params := typed.Params{
		Text:    b.text,
		ImageID: b.imageID,
		Blanks:  b.blanks,
	}

	q, err := typed.New(params)
	assert.ErrorIs(t, err, wantErr)
	return q
}
