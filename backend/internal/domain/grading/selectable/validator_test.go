package selectable_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/selectable"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	q := mustSelectableQuestion(t, []option.Option{mustSelectableOption(t, "a", true), mustSelectableOption(t, "b", false)})
	ans := mustSelectableAnswer(t)

	wrongQuestion := mustMatchingQuestion(t)
	wrongAnswer := mustMatchingAnswer(t)

	v := selectable.NewValidator()

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
