package commands

import (
	"context"
	"fmt"

	quizports "gitflic.ru/lms/backend/internal/application/ports/quiz"
	"gitflic.ru/lms/backend/internal/application/usecases/quiz/common"
	quizdomain "gitflic.ru/lms/backend/internal/domain/quiz"
)

type CreateUseCase struct {
	r         quizports.Repository
	inspector quizports.QuestionBankInspector
}

func NewCreateUseCase(r quizports.Repository, inspector quizports.QuestionBankInspector) *CreateUseCase {
	if r == nil {
		panic("quiz create usecase requires repository")
	}
	if inspector == nil {
		panic("quiz create usecase requires question bank inspector")
	}
	return &CreateUseCase{r: r, inspector: inspector}
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*Output, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}

	t, err := buildTitle(in.Title)
	if err != nil {
		return nil, err
	}

	attempts, err := buildAttempts(in.MaxAttempts)
	if err != nil {
		return nil, err
	}

	timeLimit, err := buildTimeLimit(in.TimeLimitSeconds)
	if err != nil {
		return nil, err
	}

	sources, err := buildSources(ctx, uc.inspector, in.Sources)
	if err != nil {
		return nil, err
	}

	q, err := quizdomain.New(quizdomain.Params{
		Title:            t,
		MaxAttempts:      attempts,
		TimeLimit:        timeLimit,
		ShuffleQuestions: in.ShuffleQuestions,
		Sources:          sources,
	})
	if err != nil {
		return nil, fmt.Errorf("create quiz aggregate: %w", err)
	}

	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("save quiz: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}
