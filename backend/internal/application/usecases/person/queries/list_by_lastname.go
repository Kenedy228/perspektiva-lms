package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ListByLastNameQuery struct {
	s person.QueryService
}

func NewListByLastNameQuery(s person.QueryService) *ListByLastNameQuery {
	if s == nil {
		panic("person list by last name query requires query service")
	}

	return &ListByLastNameQuery{
		s: s,
	}
}

type ListByLastnameInput struct {
	ActorRole role.Role
	LastName  string
	Limit     int
	Offset    int
}

type ListByLastNameOutput struct {
	Views []person.PersonShortView
}

func (q *ListByLastNameQuery) Execute(ctx context.Context, in ListByLastnameInput) (*ListByLastNameOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	lastName := common.NormalizeSearchText(in.LastName)
	views, err := q.s.ListByLastName(ctx, lastName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list persons by last name: %w", err)
	}

	return &ListByLastNameOutput{
		Views: views,
	}, nil
}
