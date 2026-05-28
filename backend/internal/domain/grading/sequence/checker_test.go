package sequence_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/sequence"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	selectableanswer "gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
	selectableoption "gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	opt1 := mustSequenceOption(t, "первый")
	opt2 := mustSequenceOption(t, "второй")
	opt3 := mustSequenceOption(t, "третий")
	q := mustSequenceQuestion(t, []option.Option{opt1, opt2, opt3})

	wrongQuestion := mustSelectableQuestion(t)
	wrongAnswer := mustSelectableAnswer(t)

	tests := []struct {
		name      string
		question  question.Question
		answer    question.Answer
		wantScore float64
		wantErr   error
	}{
		{name: "полностью правильный", question: q, answer: mustSequenceAnswer(t, opt1, opt2, opt3), wantScore: 1},
		{name: "неверный порядок", question: q, answer: mustSequenceAnswer(t, opt2, opt1, opt3), wantScore: 0},
		{name: "частично правильный порядок", question: q, answer: mustSequenceAnswer(t, opt1, opt3, opt2), wantScore: 0},
		{name: "недостаточное количество", question: q, answer: mustSequenceAnswer(t, opt1, opt2), wantScore: 0},
		{name: "лишние ответы", question: q, answer: mustSequenceAnswerWithIDs(t, optionID(opt1.Value()), optionID(opt2.Value()), optionID(opt3.Value()), uuid.New()), wantScore: 0},
		{name: "пустой ответ", question: q, answer: mustSequenceAnswer(t), wantScore: 0},
		{name: "неверный тип вопроса", question: wrongQuestion, answer: mustSequenceAnswer(t, opt1, opt2, opt3), wantErr: grading.ErrInvalidQuestionType},
		{name: "неверный тип ответа", question: q, answer: wrongAnswer, wantErr: grading.ErrInvalidAnswerType},
	}

	checker := sequence.New()
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

func mustBase(t *testing.T) *base.Base {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)

	b, err := base.New(ttl)
	require.NoError(t, err)
	return b
}

func mustSequenceQuestion(t *testing.T, options []option.Option) *qsequence.Question {
	t.Helper()

	q, err := qsequence.New(mustBase(t), options)
	require.NoError(t, err)
	return q
}

func mustSequenceOption(t *testing.T, value string) option.Option {
	t.Helper()

	opt, err := option.New(value)
	require.NoError(t, err)
	return opt
}

func mustSequenceAnswer(t *testing.T, options ...option.Option) answer.Answer {
	t.Helper()

	ids := make([]uuid.UUID, 0, len(options))
	for i := range options {
		ids = append(ids, optionID(options[i].Value()))
	}

	return mustSequenceAnswerWithIDs(t, ids...)
}

func mustSequenceAnswerWithIDs(t *testing.T, ids ...uuid.UUID) answer.Answer {
	t.Helper()

	optionIDs := make([]answer.OptionID, 0, len(ids))
	for i := range ids {
		id, err := answer.NewOptionID(ids[i])
		require.NoError(t, err)
		optionIDs = append(optionIDs, id)
	}

	a, err := answer.New(optionIDs)
	require.NoError(t, err)
	return a
}

func optionID(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(value))
}

func mustSelectableQuestion(t *testing.T) *qselectable.Question {
	t.Helper()

	opt1, err := selectableoption.New("верно", true)
	require.NoError(t, err)
	opt2, err := selectableoption.New("неверно", false)
	require.NoError(t, err)

	q, err := qselectable.New(mustBase(t), []selectableoption.Option{opt1, opt2})
	require.NoError(t, err)
	return q
}

func mustSelectableAnswer(t *testing.T) selectableanswer.Answer {
	t.Helper()

	a, err := selectableanswer.New(nil)
	require.NoError(t, err)
	return a
}
