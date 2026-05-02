package queries

import (
	"context"

	personqs "gitflic.ru/lms/internal/application/ports/person"
)

type ListBySnilsQuery struct {
	s personqs.QueryService
}

func NewListBySnilsQuery(s personqs.QueryService) *ListBySnilsQuery {
	return &ListBySnilsQuery{
		s: s,
	}
}

type ListBySnilsInput struct {
	Snils  string
	Limit  int
	Offset int
}

type ListBySnilsOutput struct {
	Views []personqs.PersonShortView
}

func (qs *ListBySnilsQuery) Execute(ctx context.Context, in ListBySnilsInput) (*ListBySnilsOutput, error) {
	views, err := qs.s.ListBySnils(ctx, in.Snils, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	return &ListBySnilsOutput{
		Views: views,
	}, nil
}
