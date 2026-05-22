//go:build legacy
// +build legacy

package matching_test

import (
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/question/content"
	"gitflic.ru/lms/backend/internal/domain/question/title"
)

type questionBuilder struct {
	title title.Title
	pairs []pair.Pair
}

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{}
}

func (b *questionBuilder) withTitle(t string) *questionBuilder {
	title, _ := title.New(makeContent(content.TypeText, t))
	b.title = title
	return b
}

func (b *questionBuilder) withPairs(count int) *questionBuilder {
	b.pairs = append(b.pairs, makePairs(count)...)

	return b
}

func (b *questionBuilder) build() (*matching.Question, error) {
	q, err := matching.New(b.title, b.pairs)
	return q, err
}
