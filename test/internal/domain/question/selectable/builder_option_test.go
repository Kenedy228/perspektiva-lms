package selectable_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/selectable"
	"github.com/stretchr/testify/assert"
)

func newOptionBuilder() *optionBuilder {
	return &optionBuilder{
		content:   makeContent("opt"),
		isCorrect: true,
	}
}

type optionBuilder struct {
	content   question.Content
	isCorrect bool
}

func (b *optionBuilder) withCorrect(isCorrect bool) *optionBuilder {
	b.isCorrect = isCorrect
	return b
}

func (b *optionBuilder) withContent(cType question.ContentType, val string) *optionBuilder {
	content, _ := question.NewContent(cType, val)
	b.content = content
	return b
}

func (b *optionBuilder) withContentAsText(val string) *optionBuilder {
	return b.withContent(question.ContentTypeText, val)
}

func (b *optionBuilder) build(t *testing.T, wantErr error) selectable.Option {
	t.Helper()
	params := selectable.OptionParams{
		Content:   b.content,
		IsCorrect: b.isCorrect,
	}

	item, err := selectable.NewItem(params)
	assert.ErrorIs(t, err, wantErr)

	return item
}

func (b *optionBuilder) buildNoTest() selectable.Option {
	params := selectable.OptionParams{
		Content:   b.content,
		IsCorrect: b.isCorrect,
	}

	item, _ := selectable.NewItem(params)
	return item
}

func mockOptions() []selectable.Option {
	options := make([]selectable.Option, 0, 5)

	for i := range 5 {
		options = append(options, newOptionBuilder().
			withContentAsText(fmt.Sprintf("%d", i)).
			withCorrect(true).
			buildNoTest())
	}

	return options
}
