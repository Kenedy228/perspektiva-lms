package queries

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type ListInput struct {
	ActorRole     role.Role
	AccountID     string
	TitleContains string
	Status        string
	Limit         int
	Offset        int
}

type ListOutput struct {
	Views []courseports.ShortView
}

type ListQuery struct {
	s courseports.QueryService
}

func NewListQuery(s courseports.QueryService) *ListQuery {
	if s == nil {
		panic("course list query requires query service")
	}
	return &ListQuery{s: s}
}

func (q *ListQuery) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}
	filter := courseports.Filter{
		TitleContains: common.NormalizeSearchText(in.TitleContains),
		Status:        in.Status,
		Limit:         limit,
		Offset:        offset,
	}

	switch in.ActorRole.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		views, err := q.s.ListManageable(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("list manageable courses: %w", err)
		}
		return &ListOutput{Views: views}, nil
	case role.TypeStudent:
		accountID, err := uuid.Parse(in.AccountID)
		if err != nil {
			return nil, fmt.Errorf("parse account id: %w", err)
		}
		if accountID == uuid.Nil {
			return nil, fmt.Errorf("%w: идентификатор аккаунта обязателен", common.ErrInvalidInput)
		}
		filter.AccountID = accountID
		views, err := q.s.ListVisibleForStudent(ctx, filter)
		if err != nil {
			return nil, fmt.Errorf("list student courses: %w", err)
		}
		return &ListOutput{Views: views}, nil
	default:
		return nil, fmt.Errorf("%w: роль не поддерживается", common.ErrForbidden)
	}
}
