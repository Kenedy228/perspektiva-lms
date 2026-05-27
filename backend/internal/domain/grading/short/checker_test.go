package short_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading/short"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	sequenceanswer "gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	sequenceoption "gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	"gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecker_Check(t *testing.T) {
	q := mustShortQuestion(t, "Москва")
	wrongQuestion := mustSequenceQuestion(t)
	wrongAnswer := mustSequenceAnswer(t)

	tests := []struct {
		name      string
		checker   short.Checker
		question  question.Question
		answer    question.Answer
		wantScore float64
		wantErr   error
	}{
		{name: "полностью правильный", checker: short.New(), question: q, answer: answer.New("Москва"), wantScore: 1},
		{name: "полностью неправильный", checker: short.New(), question: q, answer: answer.New("Париж"), wantScore: 0},
		{name: "пустой ответ", checker: short.New(), question: q, answer: answer.New(""), wantScore: 0},
		{name: "нормализация короткого ответа", checker: short.New(short.TrimSpace(), short.ToLower()), question: q, answer: answer.New("  москва  "), wantScore: 1},
		{name: "неверный тип вопроса", checker: short.New(), question: wrongQuestion, answer: answer.New("Москва"), wantErr: short.ErrInvalidQuestionType},
		{name: "неверный тип ответа", checker: short.New(), question: q, answer: wrongAnswer, wantErr: short.ErrInvalidAnswerType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.checker.Check(tt.question, tt.answer)
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
	checker := short.New()
	assert.True(t, checker.Supports(question.TypeShort))
	assert.False(t, checker.Supports(question.TypeSequence))
}

func mustBase(t *testing.T) *base.Base {
	t.Helper()

	ttl, err := title.New("Вопрос")
	require.NoError(t, err)

	b, err := base.New(ttl)
	require.NoError(t, err)
	return b
}

func mustShortQuestion(t *testing.T, values ...string) *qshort.Question {
	t.Helper()

	variants := make([]variant.Variant, 0, len(values))
	for i := range values {
		v, err := variant.New(values[i])
		require.NoError(t, err)
		variants = append(variants, v)
	}

	q, err := qshort.New(mustBase(t), variants)
	require.NoError(t, err)
	return q
}

func mustSequenceQuestion(t *testing.T) *qsequence.Question {
	t.Helper()

	opt1, err := sequenceoption.New("первый")
	require.NoError(t, err)
	opt2, err := sequenceoption.New("второй")
	require.NoError(t, err)

	q, err := qsequence.New(mustBase(t), []sequenceoption.Option{opt1, opt2})
	require.NoError(t, err)
	return q
}

func mustSequenceAnswer(t *testing.T) sequenceanswer.Answer {
	t.Helper()

	opt, err := sequenceanswer.NewOptionID(uuid.New())
	require.NoError(t, err)

	a, err := sequenceanswer.New([]sequenceanswer.OptionID{opt})
	require.NoError(t, err)
	return a
}
