package commands

import (
	"context"
	"fmt"

	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	"gitflic.ru/lms/backend/internal/application/usecases/quiz/common"
	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
	"gitflic.ru/lms/backend/internal/domain/quiz/limit"
	"gitflic.ru/lms/backend/internal/domain/quiz/source"
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"gitflic.ru/lms/backend/internal/domain/quiz/title"
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

func parseRequiredUUIDs(values []string, field string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for i := range values {
		id, err := parseRequiredUUID(values[i], field)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func loadQuiz(ctx context.Context, r quizports.Repository, id string) (*quizdomain.Quiz, error) {
	qID, err := parseRequiredUUID(id, "quiz id")
	if err != nil {
		return nil, err
	}

	q, err := r.FindByID(ctx, qID)
	if err != nil {
		return nil, fmt.Errorf("find quiz: %w", err)
	}

	return q, nil
}

func buildTitle(value string) (title.Title, error) {
	t, err := title.New(value)
	if err != nil {
		return title.Title{}, fmt.Errorf("create quiz title: %w", err)
	}
	return t, nil
}

func buildAttempts(count int) (limit.Attempts, error) {
	attempts, err := limit.NewAttempts(count)
	if err != nil {
		return limit.Attempts{}, fmt.Errorf("create quiz attempts limit: %w", err)
	}
	return attempts, nil
}

func buildTimeLimit(seconds int) (limit.Time, error) {
	timeLimit, err := limit.NewTime(seconds)
	if err != nil {
		return limit.Time{}, fmt.Errorf("create quiz time limit: %w", err)
	}
	return timeLimit, nil
}

func buildSources(ctx context.Context, inspector quizports.QuestionBankInspector, items []SourceInput) ([]source.Source, error) {
	sources := make([]source.Source, 0, len(items))
	for i := range items {
		s, err := buildSource(ctx, inspector, items[i])
		if err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, nil
}

func buildSource(ctx context.Context, inspector quizports.QuestionBankInspector, in SourceInput) (source.Source, error) {
	bankID, err := parseRequiredUUID(in.BankID, "question bank id")
	if err != nil {
		return source.Source{}, err
	}

	selection, questionIDs, err := buildCriteria(in)
	if err != nil {
		return source.Source{}, err
	}

	if err := ensureBankCapacity(ctx, inspector, bankID, selection.QuestionCount()); err != nil {
		return source.Source{}, err
	}

	if len(questionIDs) > 0 {
		ok, err := inspector.QuestionsBelongToBank(ctx, bankID, questionIDs)
		if err != nil {
			return source.Source{}, fmt.Errorf("check quiz source question ownership: %w", err)
		}
		if !ok {
			return source.Source{}, fmt.Errorf("%w: manual quiz source contains questions outside selected bank", common.ErrInvalidInput)
		}
	}

	s, err := source.NewSource(bankID, selection)
	if err != nil {
		return source.Source{}, fmt.Errorf("create quiz source: %w", err)
	}

	return s, nil
}

func buildCriteria(in SourceInput) (criteria.Criteria, []uuid.UUID, error) {
	switch t := criteria.Type(in.CriteriaType); t {
	case criteria.TypeRandom:
		c, err := criteria.NewRandom(in.QuestionCount)
		if err != nil {
			return nil, nil, fmt.Errorf("create random quiz source criteria: %w", err)
		}
		return c, nil, nil
	case criteria.TypeManual:
		questionIDs, err := parseRequiredUUIDs(in.QuestionIDs, "question id")
		if err != nil {
			return nil, nil, err
		}
		c, err := criteria.NewManual(questionIDs)
		if err != nil {
			return nil, nil, fmt.Errorf("create manual quiz source criteria: %w", err)
		}
		return c, questionIDs, nil
	default:
		return nil, nil, fmt.Errorf("%w: unsupported quiz source criteria type %q", common.ErrInvalidInput, in.CriteriaType)
	}
}

func ensureBankCapacity(ctx context.Context, inspector quizports.QuestionBankInspector, bankID uuid.UUID, requested int) error {
	actual, err := inspector.CountQuestionsInBank(ctx, bankID)
	if err != nil {
		return fmt.Errorf("count quiz source questions: %w", err)
	}
	if requested > actual {
		return fmt.Errorf("%w: quiz source requests %d questions, but bank contains %d", common.ErrInvalidInput, requested, actual)
	}
	return nil
}
