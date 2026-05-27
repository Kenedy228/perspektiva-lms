package sequence_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/sequence"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	opt1 := mustSequenceOption(t, "первый")
	opt2 := mustSequenceOption(t, "второй")
	q := mustSequenceQuestion(t, []option.Option{opt1, opt2})
	ans := mustSequenceAnswer(t, opt1, opt2)

	wrongQuestion := mustSelectableQuestion(t)
	wrongAnswer := mustSelectableAnswer(t)

	v := sequence.NewValidator()

	tests := []struct {
		name     string
		question question.Question
		answer   question.Answer
		wantErr  error
	}{
		{name: "ok", question: q, answer: ans},
		{name: "nil question", question: nil, answer: ans, wantErr: grading.ErrNilQuestion},
		{name: "nil answer", question: q, answer: nil, wantErr: grading.ErrNilAnswer},
		{name: "wrong question type", question: wrongQuestion, answer: ans, wantErr: grading.ErrInvalidQuestionType},
		{name: "wrong answer type", question: q, answer: wrongAnswer, wantErr: grading.ErrInvalidAnswerType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Validate(tt.question, tt.answer)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tt.wantErr))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
