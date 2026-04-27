package short_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/short"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{
		text:     "text",
		imageID:  uuid.Nil,
		variants: []question.Content{},
	}
}

type questionBuilder struct {
	text     string
	imageID  uuid.UUID
	variants []question.Content
}

func (b *questionBuilder) withText(s string) *questionBuilder {
	b.text = s
	return b
}

func (b *questionBuilder) withImage(id uuid.UUID) *questionBuilder {
	b.imageID = id
	return b
}

func (b *questionBuilder) withVariant(val string) *questionBuilder {
	item := makeContent(val)
	b.variants = append(b.variants, item)
	return b
}

func (b *questionBuilder) build(t *testing.T, wantErr error) question.Question {
	params := short.Params{
		Text:     b.text,
		ImageID:  b.imageID,
		Variants: b.variants,
	}

	q, err := short.New(params)
	assert.ErrorIs(t, err, wantErr)
	return q
}
