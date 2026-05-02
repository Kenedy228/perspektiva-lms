package matching_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question/content"
	"gitflic.ru/lms/internal/domain/question/matching"
	"gitflic.ru/lms/internal/domain/question/matching/answer"
	"gitflic.ru/lms/internal/domain/question/matching/pair"
	"gitflic.ru/lms/internal/domain/question/title"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func makeAnswer(pairs map[uuid.UUID]uuid.UUID) answer.Answer {
	var ansPairs []answer.AnswerPair
	for promptID, matchID := range pairs {
		ansPairs = append(ansPairs, answer.AnswerPair{
			PromptID: promptID,
			MatchID:  matchID,
		})
	}
	return answer.New(ansPairs)
}

type questionBuilder struct {
	titleStr string
	pairs    []pair.Pair
}

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{}
}

func (b *questionBuilder) withTitle(s string) *questionBuilder {
	b.titleStr = s
	return b
}

func (b *questionBuilder) withPair(promptVal, matchVal string) *questionBuilder {
	promptContent, _ := content.New(content.TypeText, promptVal)
	matchContent, _ := content.New(content.TypeText, matchVal)

	p, _ := pair.New(promptContent, matchContent)
	b.pairs = append(b.pairs, p)
	return b
}

func (b *questionBuilder) build(t *testing.T) *matching.Question {
	t.Helper()

	tTitle, err := title.New(makeContent(b.titleStr))
	require.NoError(t, err)

	q, err := matching.New(tTitle, b.pairs)
	require.NoError(t, err)

	return q
}

func makeContent(value string) content.Content {
	c, _ := content.New(content.TypeText, value)
	return c
}
