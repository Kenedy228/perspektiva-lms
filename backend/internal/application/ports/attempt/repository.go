package attempt

import (
	"context"
	"time"

	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*attemptdomain.Attempt, error)
	Save(ctx context.Context, a *attemptdomain.Attempt) error
	CountByEnrollmentAndQuiz(ctx context.Context, enrollmentID, quizID uuid.UUID) (int, error)
}

type EnrollmentPolicy interface {
	CanStartQuiz(ctx context.Context, accountID, enrollmentID, quizID uuid.UUID, at time.Time) (bool, error)
}

type QuestionSetProvider interface {
	FindQuestionsByIDs(ctx context.Context, bankID uuid.UUID, questionIDs []uuid.UUID) ([]question.Question, error)
	SelectRandomQuestions(ctx context.Context, bankID uuid.UUID, count int) ([]question.Question, error)
}

type QuestionShuffler interface {
	ShuffleQuestions(questions []question.Question) []question.Question
}
