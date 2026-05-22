package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/person/profile/snils"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ListBySnilsQuery struct {
	s person.QueryService
}

func NewListBySnilsQuery(s person.QueryService) *ListBySnilsQuery {
	if s == nil {
		panic("person list by snils query requires query service")
	}

	return &ListBySnilsQuery{
		s: s,
	}
}

type ListBySnilsInput struct {
	ActorRole role.Role
	Snils     string
	Limit     int
	Offset    int
}

type ListBySnilsOutput struct {
	Views []person.PersonShortView
}

func (q *ListBySnilsQuery) Execute(ctx context.Context, in ListBySnilsInput) (*ListBySnilsOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	s, err := snils.New(common.NormalizeSNILSSearch(in.Snils))
	if err != nil {
		return nil, fmt.Errorf("list persons by snils value: %w", err)
	}

	views, err := q.s.ListBySnils(ctx, s.Value(), limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list persons by snils: %w", err)
	}

	return &ListBySnilsOutput{
		Views: views,
	}, nil
}
