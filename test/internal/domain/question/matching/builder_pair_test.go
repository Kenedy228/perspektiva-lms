package matching_test

import (
	"fmt"
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/matching"
	"github.com/stretchr/testify/assert"
)

func newPairBuilder() *pairBuilder {
	return &pairBuilder{
		prompt:  defaultPairPrompt(),
		content: defaultPairOption(),
	}
}

type pairBuilder struct {
	prompt  string
	content question.Content
}

func (b *pairBuilder) withPrompt(prompt string) *pairBuilder {
	b.prompt = prompt
	return b
}

func (b *pairBuilder) withContent(cType question.ContentType, s string) *pairBuilder {
	content, err := question.NewContent(cType, s)
	if err != nil {
		panic(err)
	}

	b.content = content
	return b
}

func (b *pairBuilder) withContentText(s string) *pairBuilder {
	b.content = makeContent(s)
	return b
}

func (b *pairBuilder) build(t *testing.T, wantErr error) matching.Pair {
	t.Helper()

	params := matching.PairParams{
		Prompt:  b.prompt,
		Content: b.content,
	}
	pair, err := matching.NewPair(params)
	assert.ErrorIs(t, err, wantErr)

	return pair
}

func (b *pairBuilder) buildNoTest() matching.Pair {
	params := matching.PairParams{
		Prompt:  b.prompt,
		Content: b.content,
	}

	pair, _ := matching.NewPair(params)
	return pair
}

func mockPairs() []matching.Pair {
	pairs := make([]matching.Pair, 0, 5)

	for i := range 5 {
		v := fmt.Sprintf("%d", i)
		p := newPairBuilder().withPrompt(v).withContentText(v).buildNoTest()

		pairs = append(pairs, p)
	}

	return pairs
}
