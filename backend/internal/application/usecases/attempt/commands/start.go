package commands

import (
	"context"
	"fmt"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	"gitflic.ru/lms/backend/internal/application/usecases/attempt/common"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
)

type StartUseCase struct {
	attempts   attemptports.Repository
	quizzes    quizports.Repository
	enrollment attemptports.EnrollmentPolicy
	questions  attemptports.QuestionSetProvider
	shuffler   attemptports.QuestionShuffler
}

func NewStartUseCase(
	attempts attemptports.Repository,
	quizzes quizports.Repository,
	enrollment attemptports.EnrollmentPolicy,
	questions attemptports.QuestionSetProvider,
	shuffler attemptports.QuestionShuffler,
) *StartUseCase {
	if attempts == nil {
		panic("attempt start usecase requires attempt repository")
	}
	if quizzes == nil {
		panic("attempt start usecase requires quiz repository")
	}
	if enrollment == nil {
		panic("attempt start usecase requires enrollment policy")
	}
	if questions == nil {
		panic("attempt start usecase requires question set provider")
	}
	if shuffler == nil {
		panic("attempt start usecase requires question shuffler")
	}
	return &StartUseCase{
		attempts:   attempts,
		quizzes:    quizzes,
		enrollment: enrollment,
		questions:  questions,
		shuffler:   shuffler,
	}
}

func (uc *StartUseCase) Execute(ctx context.Context, in StartInput) (*Output, error) {
	if err := common.RequireStudent(in.ActorRole); err != nil {
		return nil, err
	}

	accountID, err := parseRequiredUUID(in.AccountID, "account id")
	if err != nil {
		return nil, err
	}
	enrollmentID, err := parseRequiredUUID(in.EnrollmentID, "enrollment id")
	if err != nil {
		return nil, err
	}
	quizID, err := parseRequiredUUID(in.QuizID, "quiz id")
	if err != nil {
		return nil, err
	}

	ok, err := uc.enrollment.CanStartQuiz(ctx, accountID, enrollmentID, quizID, in.Now)
	if err != nil {
		return nil, fmt.Errorf("check quiz enrollment: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("%w: student is not enrolled for this quiz", common.ErrForbidden)
	}

	q, err := loadQuiz(ctx, uc.quizzes, quizID)
	if err != nil {
		return nil, err
	}

	started, err := uc.attempts.CountByEnrollmentAndQuiz(ctx, enrollmentID, quizID)
	if err != nil {
		return nil, fmt.Errorf("count quiz attempts: %w", err)
	}
	if !q.Attempts().IsInfinite() && started >= q.Attempts().Count() {
		return nil, fmt.Errorf("%w: quiz attempts limit reached", common.ErrLimitReached)
	}

	questions, err := materializeQuestions(ctx, uc.questions, q.Sources())
	if err != nil {
		return nil, err
	}
	if q.ShuffleQuestions() {
		questions = uc.shuffler.ShuffleQuestions(questions)
	}

	a, err := attemptdomain.New(attemptdomain.Params{
		EnrollmentID: enrollmentID,
		QuizID:       quizID,
		TimeLimit:    q.Time(),
		Questions:    questions,
	}, in.Now)
	if err != nil {
		return nil, fmt.Errorf("create attempt aggregate: %w", err)
	}

	if err := saveAttempt(ctx, uc.attempts, a); err != nil {
		return nil, err
	}

	return &Output{ID: a.ID().String()}, nil
}
