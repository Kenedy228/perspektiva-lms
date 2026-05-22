package matching_test

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_RejectsDuplicatePairs(t *testing.T) {
	t1, err := title.New("Вопрос")
	require.NoError(t, err)
	p := mustPair(t, "a", "1")

	_, err = matching.New(t1, []pair.Pair{p, p})
	assert.ErrorIs(t, err, matching.ErrInvalid)
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
