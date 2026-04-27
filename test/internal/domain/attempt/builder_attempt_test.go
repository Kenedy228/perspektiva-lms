package attempt_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/attempt"
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/shared/limit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type attemptBuilder struct {
	questions    []question.Question
	enrollmentID uuid.UUID
	quizID       uuid.UUID
	timeLimit    limit.Limit
}

func newAttemptBuilder() *attemptBuilder {
	return &attemptBuilder{
		questions:    []question.Question{},
		enrollmentID: uuid.Nil,
		quizID:       uuid.Nil,
		timeLimit:    limit.Limit{},
	}
}

func (b *attemptBuilder) withMockQuestions(size int) *attemptBuilder {
	questions := make([]question.Question, 0, size)

	for range size {
		q := new(mockQuestion)
		q.On("Clone").Return(q)
		questions = append(questions, q)
	}

	b.questions = questions
	return b
}

func (b *attemptBuilder) withQuestions(questions []question.Question) *attemptBuilder {
	b.questions = questions
	return b
}

func (b *attemptBuilder) withEnrollment(id uuid.UUID) *attemptBuilder {
	b.enrollmentID = id
	return b
}

func (b *attemptBuilder) withQuiz(id uuid.UUID) *attemptBuilder {
	b.quizID = id
	return b
}

func (b *attemptBuilder) withTimeLimit(val int) *attemptBuilder {
	l, _ := limit.New(val)
	b.timeLimit = l
	return b
}

func (b *attemptBuilder) build(t *testing.T, wantErr error) *attempt.Attempt {
	t.Helper()

	params := attempt.Params{
		EnrollmentID: b.enrollmentID,
		QuizID:       b.quizID,
		TimeLimit:    b.timeLimit,
		Questions:    b.questions,
	}

	a, err := attempt.New(params)
	assert.ErrorIs(t, err, wantErr)

	return a
}
