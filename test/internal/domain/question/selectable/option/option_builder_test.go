//go:build legacy
// +build legacy

package option_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"gitflic.ru/lms/backend/internal/domain/question/content"
)

type optionBuilder struct {
	c         content.Content
	isCorrect bool
}

func newOptionBuilder() *optionBuilder {
	return &optionBuilder{}
}

func (b *optionBuilder) withContent(cType content.Type, s string) *optionBuilder {
	c, _ := content.New(cType, s)
	b.c = c
	return b
}

func (b *optionBuilder) withIsCorrect(isCorrect bool) *optionBuilder {
	b.isCorrect = isCorrect
	return b
}

func (b *optionBuilder) build() (option.Option, error) {
	return option.New(b.c, b.isCorrect)
}
