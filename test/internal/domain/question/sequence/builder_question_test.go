package sequence_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/sequence"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{
		text:     "text",
		imageID:  uuid.Nil,
		elements: []sequence.Element{},
	}
}

type questionBuilder struct {
	text     string
	imageID  uuid.UUID
	elements []sequence.Element
}

func (b *questionBuilder) withText(s string) *questionBuilder {
	b.text = s
	return b
}

func (b *questionBuilder) withImage(id uuid.UUID) *questionBuilder {
	b.imageID = id
	return b
}

func (b *questionBuilder) withElement(cType question.ContentType, val string) *questionBuilder {
	item := newElementBuilder().withContent(cType, val).
		buildNoTest()

	b.elements = append(b.elements, item)
	return b
}

func (b *questionBuilder) withElementAsText(val string) *questionBuilder {
	item := newElementBuilder().withContentAsText(val).
		buildNoTest()

	b.elements = append(b.elements, item)
	return b
}

func (b *questionBuilder) build(t *testing.T, wantErr error) question.Question {
	params := sequence.Params{
		Text:     b.text,
		ImageID:  b.imageID,
		Elements: b.elements,
	}

	q, err := sequence.New(params)
	assert.ErrorIs(t, err, wantErr)
	return q
}
