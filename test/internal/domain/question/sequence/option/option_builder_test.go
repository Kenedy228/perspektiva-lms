package option_test

import (
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/sequence/option"
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

func (b *optionBuilder) build() (option.Option, error) {
	return option.New(b.c)
}
