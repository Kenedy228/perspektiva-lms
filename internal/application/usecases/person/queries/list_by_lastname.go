package queries

import (
	"context"

	personqs "gitflic.ru/lms/internal/application/ports/person"
)

type ListByLastNameQuery struct {
	s personqs.QueryService
}

func NewListByLastNameQuery(s personqs.QueryService) *ListByLastNameQuery {
	return &ListByLastNameQuery{
		s: s,
	}
}

type ListByLastnameInput struct {
	LastName string
	Limit    int
	Offset   int
}

type ListByLastNameOutput struct {
	Views []personqs.PersonShortView
}

func (qs *ListByLastNameQuery) Execute(ctx context.Context, in ListByLastnameInput) (*ListByLastNameOutput, error) {
	views, err := qs.s.ListByLastName(ctx, in.LastName, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	return &ListByLastNameOutput{
		Views: views,
	}, nil
}
