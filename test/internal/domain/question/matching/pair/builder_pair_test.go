//go:build legacy
// +build legacy

package pair_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/question/content"
)

type pairBuilder struct {
	prompt content.Content
	match  content.Content
}

func newPairBuilder() *pairBuilder {
	return &pairBuilder{}
}

func (b *pairBuilder) withPrompt(cType content.Type, value string) *pairBuilder {
	b.prompt = makeContent(cType, value)
	return b
}

func (b *pairBuilder) withMatch(cType content.Type, value string) *pairBuilder {
	b.match = makeContent(cType, value)
	return b
}

func (b *pairBuilder) build() (pair.Pair, error) {
	return pair.New(b.prompt, b.match)
}
