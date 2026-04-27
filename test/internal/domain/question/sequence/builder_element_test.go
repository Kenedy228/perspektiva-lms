package sequence_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/sequence"
	"github.com/stretchr/testify/assert"
)

func newElementBuilder() *elementBuilder {
	return &elementBuilder{
		content:   makeContent("opt"),
		isCorrect: true,
	}
}

type elementBuilder struct {
	content   question.Content
	isCorrect bool
}

func (b *elementBuilder) withContent(cType question.ContentType, val string) *elementBuilder {
	opt, _ := question.NewContent(cType, val)
	b.content = opt
	return b
}

func (b *elementBuilder) withContentAsText(val string) *elementBuilder {
	return b.withContent(question.ContentTypeText, val)
}

func (b *elementBuilder) build(t *testing.T, wantErr error) sequence.Element {
	t.Helper()
	params := sequence.ElementParams{
		Content: b.content,
	}

	item, err := sequence.NewElement(params)
	assert.ErrorIs(t, err, wantErr)

	return item
}

func (b *elementBuilder) buildNoTest() sequence.Element {
	params := sequence.ElementParams{
		Content: b.content,
	}

	item, _ := sequence.NewElement(params)
	return item
}

func mockElements() []sequence.Element {
	elements := make([]sequence.Element, 0, 5)
	for i := range 5 {
		elements = append(elements, newElementBuilder().
			withContentAsText(fmt.Sprintf("%d", i)).
			buildNoTest())
	}

	return elements
}
