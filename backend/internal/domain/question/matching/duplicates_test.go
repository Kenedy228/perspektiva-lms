package matching_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	basetitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"github.com/stretchr/testify/require"
)

func TestNew_AllowsDuplicatePairs(t *testing.T) {
	t1, err := basetitle.New("Вопрос")
	require.NoError(t, err)

	b, err := base.New(t1)
	require.NoError(t, err)

	p := mustPair(t, "a", "1")

	q, err := matching.New(b, []pair.Pair{p, p})
	require.NoError(t, err)
	require.Len(t, q.Pairs(), 2)
}

func mustPair(t *testing.T, promptValue, matchValue string) pair.Pair {
	t.Helper()

	prompt, err := pair.NewPrompt(promptValue)
	require.NoError(t, err)
	match, err := pair.NewMatch(matchValue)
	require.NoError(t, err)

	p, err := pair.New(prompt, match)
	require.NoError(t, err)
	return p
}
