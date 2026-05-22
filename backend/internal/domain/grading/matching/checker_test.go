package matching_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/matching"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	t1, err := title.New("Вопрос")
	require.NoError(t, err)
	p1 := mustPair(t, "a", "1")
	p2 := mustPair(t, "b", "2")
	q, err := qmatching.New(t1, []pair.Pair{p1, p2})
	require.NoError(t, err)

	prompt1, err := answer.NewPromptID(p1.PromptID())
	require.NoError(t, err)
	match1, err := answer.NewMatchID(p1.MatchID())
	require.NoError(t, err)
	prompt2, err := answer.NewPromptID(p2.PromptID())
	require.NoError(t, err)
	match2, err := answer.NewMatchID(p2.MatchID())
	require.NoError(t, err)
	ans, err := answer.New([]answer.Pair{
		{PromptID: prompt1, MatchID: match1},
		{PromptID: prompt2, MatchID: match2},
	})
	require.NoError(t, err)

	got, err := matching.New().Check(q, ans)
	require.NoError(t, err)
	assert.Equal(t, 1.0, got.Value())
}

func mustPair(t *testing.T, promptValue, matchValue string) pair.Pair {
	t.Helper()

	promptText, err := text.New(promptValue)
	require.NoError(t, err)
	matchText, err := text.New(matchValue)
	require.NoError(t, err)

	prompt, err := pair.NewPrompt(promptText)
	require.NoError(t, err)
	match, err := pair.NewMatch(matchText)
	require.NoError(t, err)

	p, err := pair.New(prompt, match)
	require.NoError(t, err)
	return p
}
