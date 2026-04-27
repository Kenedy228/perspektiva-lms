package typed_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/typed"
	"github.com/stretchr/testify/assert"
)

func newBlankBuilder() *blankBuilder {
	return &blankBuilder{
		placeholder: "{{placeholder}}",
		variants:    []question.Content{},
	}
}

type blankBuilder struct {
	placeholder string
	variants    []question.Content
}

func (b *blankBuilder) withPlaceholder(s string) *blankBuilder {
	b.placeholder = fmt.Sprintf("{{%s}}", s)
	return b
}

func (b *blankBuilder) withInvalidPlaceholder(s string) *blankBuilder {
	b.placeholder = fmt.Sprintf("%s", s)
	return b
}

func (b *blankBuilder) withVariant(s string) *blankBuilder {
	b.variants = append(b.variants, makeContent(s))
	return b
}

func (b *blankBuilder) withInvalidVariant() *blankBuilder {
	opt, _ := question.NewContent(question.ContentTypeImage, "image")
	b.variants = append(b.variants, opt)
	return b
}

func (b *blankBuilder) build(t *testing.T, wantErr error) typed.Blank {
	params := typed.BlankParams{
		Placeholder: b.placeholder,
		Variants:    b.variants,
	}

	blank, err := typed.NewBlank(params)
	assert.ErrorIs(t, err, wantErr)

	return blank
}

func (b *blankBuilder) buildNoTest() typed.Blank {
	params := typed.BlankParams{
		Placeholder: b.placeholder,
		Variants:    b.variants,
	}

	blank, _ := typed.NewBlank(params)
	return blank
}
