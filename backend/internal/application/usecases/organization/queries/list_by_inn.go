package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ListByINNQuery struct {
	s organization.QueryService
}

func NewListByINNQuery(s organization.QueryService) *ListByINNQuery {
	if s == nil {
		panic("organization list by inn query requires query service")
	}

	return &ListByINNQuery{
		s: s,
	}
}

type ListByINNInput struct {
	ActorRole role.Role
	INN       string
	Limit     int
	Offset    int
}

type ListByINNOutput struct {
	Views []organization.OrganizationShortView
}

func (q *ListByINNQuery) Execute(ctx context.Context, in ListByINNInput) (*ListByINNOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	search := common.NormalizeINNSearch(in.INN)
	if search == "" {
		return nil, fmt.Errorf("%w: inn search cannot be empty", common.ErrInvalidInput)
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	views, err := q.s.ListByINN(ctx, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list organizations by inn: %w", err)
	}

	return &ListByINNOutput{
		Views: views,
	}, nil
}
