package matching_test

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/matching"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/matching/answer"
	shortanswer "gitflic.ru/lms/backend/internal/domain/question/short/answer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidator_Validate(t *testing.T) {
	p1 := mustPair(t, "a", "1")
	p2 := mustPair(t, "b", "2")
	q := mustMatchingQuestion(t, p1, p2)
	ans := mustMatchingAnswer(t, []answer.Pair{{PromptID: p1.PromptID(), MatchID: p1.MatchID()}})

	wrongQuestion := mustShortQuestion(t)
	wrongAnswer := shortanswer.New("text")

	v := matching.NewValidator()

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
