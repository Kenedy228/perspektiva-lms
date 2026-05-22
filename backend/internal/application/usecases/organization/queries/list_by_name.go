package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ListByNameQuery struct {
	s organization.QueryService
}

func NewListByNameQuery(s organization.QueryService) *ListByNameQuery {
	if s == nil {
		panic("organization list by name query requires query service")
	}

	return &ListByNameQuery{
		s: s,
	}
}

type ListByNameInput struct {
	ActorRole role.Role
	Name      string
	Limit     int
	Offset    int
}

type ListByNameOutput struct {
	Views []organization.OrganizationShortView
}

func (q *ListByNameQuery) Execute(ctx context.Context, in ListByNameInput) (*ListByNameOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	search := common.NormalizeSearchText(in.Name)
	if search == "" {
		return nil, fmt.Errorf("%w: name search cannot be empty", common.ErrInvalidInput)
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	views, err := q.s.ListByName(ctx, search, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list organizations by name: %w", err)
	}

	return &ListByNameOutput{
		Views: views,
	}, nil
}
