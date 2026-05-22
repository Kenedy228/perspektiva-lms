package queries

import (
	"context"
	"fmt"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type ListQuery struct {
	s bankports.QueryService
}

func NewListQuery(s bankports.QueryService) *ListQuery {
	if s == nil {
		panic("bank list query requires query service")
	}
	return &ListQuery{s: s}
}

type ListInput struct {
	ActorRole     role.Role
	TitleContains string
	QuestionID    string
	MinQuestions  int
	MaxQuestions  int
	Limit         int
	Offset        int
}

type ListOutput struct {
	Views []bankports.ShortView
}

func (q *ListQuery) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	if in.MinQuestions < 0 {
		return nil, fmt.Errorf("%w: min questions cannot be negative", common.ErrInvalidInput)
	}
	if in.MaxQuestions < 0 {
		return nil, fmt.Errorf("%w: max questions cannot be negative", common.ErrInvalidInput)
	}
	if in.MaxQuestions > 0 && in.MinQuestions > in.MaxQuestions {
		return nil, fmt.Errorf("%w: min questions cannot exceed max questions", common.ErrInvalidInput)
	}

	var questionID uuid.UUID
	if in.QuestionID != "" {
		questionID, err = uuid.Parse(in.QuestionID)
		if err != nil {
			return nil, fmt.Errorf("parse question id: %w", err)
		}
		if questionID == uuid.Nil {
			return nil, fmt.Errorf("%w: question id is required", common.ErrInvalidInput)
		}
	}

	views, err := q.s.List(ctx, bankports.Filter{
		TitleContains: common.NormalizeSearchText(in.TitleContains),
		QuestionID:    questionID,
		MinQuestions:  in.MinQuestions,
		MaxQuestions:  in.MaxQuestions,
		Limit:         limit,
		Offset:        offset,
	})
	if err != nil {
		return nil, fmt.Errorf("list banks: %w", err)
	}

	return &ListOutput{Views: views}, nil
}
