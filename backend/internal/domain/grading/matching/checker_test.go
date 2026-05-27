package matching_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/matching"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	p1 := mustPair(t, "a", "1")
	p2 := mustPair(t, "b", "2")
	q := mustMatchingQuestion(t, p1, p2)

	wrongQuestion := mustShortQuestion(t)
	wrongAnswer := shortanswer.New("text")

	tests := []struct {
		name      string
		question  question.Question
		answer    question.Answer
		wantScore float64
		wantErr   error
	}{
		{name: "полностью правильный", question: q, answer: mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p1.MatchID()}, {PromptID: p2.PromptID(), MatchID: p2.MatchID()}}), wantScore: 1},
		{name: "полностью неправильный", question: q, answer: mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p2.MatchID()}, {PromptID: p2.PromptID(), MatchID: p1.MatchID()}}), wantScore: 0},
		{name: "частично правильный", question: q, answer: mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p1.MatchID()}, {PromptID: p2.PromptID(), MatchID: uuid.New()}}), wantScore: 0},
		{name: "лишние пары", question: q, answer: mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p1.MatchID()}, {PromptID: p2.PromptID(), MatchID: p2.MatchID()}, {PromptID: uuid.New(), MatchID: uuid.New()}}), wantScore: 0},
		{name: "недостаточное количество пар", question: q, answer: mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p1.MatchID()}}), wantScore: 0},
		{name: "пустой ответ", question: q, answer: mustMatchingAnswer(t, nil), wantScore: 0},
		{name: "неверный тип вопроса", question: wrongQuestion, answer: mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p1.MatchID()}}), wantErr: grading.ErrInvalidQuestionType},
		{name: "неверный тип ответа", question: q, answer: wrongAnswer, wantErr: grading.ErrInvalidAnswerType},
	}

	checker := matching.New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checker.Check(tt.question, tt.answer)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantScore, got.Value())
		})
	}
}

func TestChecker_Supports(t *testing.T) {
	checker := matching.New()
	assert.True(t, checker.Supports(question.TypeMatching))
	assert.False(t, checker.Supports(question.TypeSelectable))
}

func mustBase(t *testing.T) *base.Base {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)

	b, err := base.New(ttl)
	require.NoError(t, err)
	return b
}

func mustMatchingQuestion(t *testing.T, pairs ...pair.Pair) *qmatching.Question {
	t.Helper()

	q, err := qmatching.New(mustBase(t), pairs)
	require.NoError(t, err)
	return q
}

func mustMatchingAnswer(t *testing.T, pairs []answer.Pair) answer.Answer {
	t.Helper()

	a, err := answer.New(pairs)
	require.NoError(t, err)
	return a
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

func mustShortQuestion(t *testing.T) *qshort.Question {
	t.Helper()

	v, err := variant.New("ответ")
	require.NoError(t, err)

	q, err := qshort.New(mustBase(t), []variant.Variant{v})
	require.NoError(t, err)
	return q
}
