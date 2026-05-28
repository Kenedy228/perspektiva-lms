package queries

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ListQuery struct {
	s accountports.QueryService
}

func NewListQuery(s accountports.QueryService) *ListQuery {
	if s == nil {
		panic("account list query requires query service")
	}

	return &ListQuery{s: s}
}

type ListInput struct {
	ActorRole role.Role
	Role      role.Type
	Status    accountdomain.Status
	Login     string
	Limit     int
	Offset    int
}

type ListOutput struct {
	Views []accountports.AccountView
}

func (q *ListQuery) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	if in.Role != "" && !in.Role.IsValid() {
		return nil, fmt.Errorf("%w: role filter is invalid", common.ErrInvalidInput)
	}

	if in.Status != "" && !in.Status.IsValid() {
		return nil, fmt.Errorf("%w: status filter is invalid", common.ErrInvalidInput)
	}

	views, err := q.s.List(ctx, accountports.ListFilter{
		Role:   in.Role,
		Status: in.Status,
		Login:  common.NormalizeSearchText(in.Login),
	}, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}

	return &ListOutput{Views: views}, nil
}
