package matching_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/matching"
	"gitflic.ru/lms/internal/domain/question/matching/pair"
	"github.com/stretchr/testify/require"
)

func makePairs(count int) []pair.Pair {
	pairs := make([]pair.Pair, 0, count)

	for range count {
		pair, _ := pair.New(
			makeContent(content.TypeText, "text"),
			makeContent(content.TypeText, "text2"),
		)

		pairs = append(pairs, pair)
	}

	return pairs
}

func makeContent(cType content.Type, s string) content.Content {
	c, _ := content.New(cType, s)
	return c
}

func castQuestion(t *testing.T, q question.Question) *matching.Question {
	cast, ok := q.(*matching.Question)
	require.True(t, ok)
	return cast
}
