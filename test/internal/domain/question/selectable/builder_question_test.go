package selectable_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/selectable"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{
		text:    "text",
		imageID: uuid.Nil,
		options: []selectable.Option{},
	}
}

type questionBuilder struct {
	text    string
	imageID uuid.UUID
	options []selectable.Option
}

func (b *questionBuilder) withText(s string) *questionBuilder {
	b.text = s
	return b
}

func (b *questionBuilder) withImage(id uuid.UUID) *questionBuilder {
	b.imageID = id
	return b
}

func (b *questionBuilder) withOption(cType question.ContentType, val string, isCorrect bool) *questionBuilder {
	option := newOptionBuilder().withContent(cType, val).
		withCorrect(isCorrect).
		buildNoTest()

	b.options = append(b.options, option)
	return b
}

func (b *questionBuilder) withOptionAsText(val string, isCorrect bool) *questionBuilder {
	option := newOptionBuilder().withContentAsText(val).
		withCorrect(isCorrect).
		buildNoTest()

	b.options = append(b.options, option)
	return b
}

func (b *questionBuilder) build(t *testing.T, wantErr error) question.Question {
	params := selectable.Params{
		Text:    b.text,
		ImageID: b.imageID,
		Options: b.options,
	}

	q, err := selectable.New(params)
	assert.ErrorIs(t, err, wantErr)
	return q
}
