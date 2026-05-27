package selectable_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/selectable"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	matchinganswer "gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	opt1 := mustSelectableOption(t, "верно 1", true)
	opt2 := mustSelectableOption(t, "верно 2", true)
	opt3 := mustSelectableOption(t, "неверно", false)
	q := mustSelectableQuestion(t, []option.Option{opt1, opt2, opt3})

	wrongQuestion := mustMatchingQuestion(t)
	wrongAnswer := mustMatchingAnswer(t)

	tests := []struct {
		name      string
		question  question.Question
		answer    question.Answer
		wantScore float64
		wantErr   error
	}{
		{name: "полностью правильный", question: q, answer: mustSelectableAnswer(t, opt1.ID(), opt2.ID()), wantScore: 1},
		{name: "полностью неправильный", question: q, answer: mustSelectableAnswer(t, opt3.ID()), wantScore: 0},
		{name: "частично правильный", question: q, answer: mustSelectableAnswer(t, opt1.ID()), wantScore: 0},
		{name: "лишний вариант", question: q, answer: mustSelectableAnswer(t, opt1.ID(), opt2.ID(), opt3.ID()), wantScore: 0},
		{name: "пустой ответ", question: q, answer: mustSelectableAnswer(t), wantScore: 0},
		{name: "неверный тип вопроса", question: wrongQuestion, answer: mustSelectableAnswer(t, opt1.ID()), wantErr: selectable.ErrInvalidQuestionType},
		{name: "неверный тип ответа", question: q, answer: wrongAnswer, wantErr: selectable.ErrInvalidAnswerType},
	}

	checker := selectable.New()
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
	checker := selectable.New()
	assert.True(t, checker.Supports(question.TypeSelectable))
	assert.False(t, checker.Supports(question.TypeMatching))
}

func mustBase(t *testing.T) *base.Base {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)

	b, err := base.New(ttl)
	require.NoError(t, err)
	return b
}

func mustSelectableQuestion(t *testing.T, options []option.Option) *qselectable.Question {
	t.Helper()

	q, err := qselectable.New(mustBase(t), options)
	require.NoError(t, err)
	return q
}

func mustSelectableOption(t *testing.T, value string, correct bool) option.Option {
	t.Helper()

	opt, err := option.New(value, correct)
	require.NoError(t, err)
	return opt
}

func mustSelectableAnswer(t *testing.T, ids ...uuid.UUID) answer.Answer {
	t.Helper()

	a, err := answer.New(ids)
	require.NoError(t, err)
	return a
}

func mustMatchingQuestion(t *testing.T) *qmatching.Question {
	t.Helper()

	p1 := mustPair(t, "a", "1")
	p2 := mustPair(t, "b", "2")
	q, err := qmatching.New(mustBase(t), []pair.Pair{p1, p2})
	require.NoError(t, err)
	return q
}

func mustMatchingAnswer(t *testing.T) matchinganswer.Answer {
	t.Helper()

	a, err := matchinganswer.New([]matchinganswer.Pair{{PromptID: uuid.New(), MatchID: uuid.New()}})
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
