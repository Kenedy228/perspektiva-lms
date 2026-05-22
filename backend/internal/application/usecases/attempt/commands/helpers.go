package commands

import (
	"context"
	"fmt"

	attemptports "gitflic.ru/lms/backend/internal/application/ports/attempt"
	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	"gitflic.ru/lms/backend/internal/application/usecases/attempt/common"
	attemptdomain "gitflic.ru/lms/backend/internal/domain/attempt"
	"gitflic.ru/lms/backend/internal/domain/question"
	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
)

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: %s is required", common.ErrInvalidInput, field)
	}
	return id, nil
}

func loadAttempt(ctx context.Context, r attemptports.Repository, id string) (*attemptdomain.Attempt, error) {
	attemptID, err := parseRequiredUUID(id, "attempt id")
	if err != nil {
		return nil, err
	}

	a, err := r.FindByID(ctx, attemptID)
	if err != nil {
		return nil, fmt.Errorf("find attempt: %w", err)
	}

	return a, nil
}

func loadQuiz(ctx context.Context, r quizports.Repository, id uuid.UUID) (*quizdomain.Quiz, error) {
	q, err := r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find quiz: %w", err)
	}
	return q, nil
}

func saveAttempt(ctx context.Context, r attemptports.Repository, a *attemptdomain.Attempt) error {
	if err := r.Save(ctx, a); err != nil {
		return fmt.Errorf("save attempt: %w", err)
	}
	return nil
}

func materializeQuestions(ctx context.Context, provider attemptports.QuestionSetProvider, sources []source.Source) ([]question.Question, error) {
	questions := make([]question.Question, 0)
	for i := range sources {
		selected, err := materializeSource(ctx, provider, sources[i])
		if err != nil {
			return nil, err
		}
		questions = append(questions, selected...)
	}
	return questions, nil
}

func materializeSource(ctx context.Context, provider attemptports.QuestionSetProvider, s source.Source) ([]question.Question, error) {
	switch c := s.Criteria().(type) {
	case criteria.Manual:
		questions, err := provider.FindQuestionsByIDs(ctx, s.BankID(), c.QuestionIDs())
		if err != nil {
			return nil, fmt.Errorf("load manual quiz source questions: %w", err)
		}
		if len(questions) != c.QuestionCount() {
			return nil, fmt.Errorf("%w: manual quiz source returned %d questions, expected %d", common.ErrInvalidInput, len(questions), c.QuestionCount())
		}
		return questions, nil
	case criteria.Random:
		questions, err := provider.SelectRandomQuestions(ctx, s.BankID(), c.QuestionCount())
		if err != nil {
			return nil, fmt.Errorf("select random quiz source questions: %w", err)
		}
		if len(questions) != c.QuestionCount() {
			return nil, fmt.Errorf("%w: random quiz source returned %d questions, expected %d", common.ErrInvalidInput, len(questions), c.QuestionCount())
		}
		return questions, nil
	default:
		return nil, fmt.Errorf("%w: unsupported quiz source criteria", common.ErrInvalidInput)
	}
}
