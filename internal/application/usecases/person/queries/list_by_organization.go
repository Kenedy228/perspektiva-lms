package queries

import (
	"context"

	personqs "gitflic.ru/lms/internal/application/ports/person"
	"github.com/google/uuid"
)

type ListByOrganizationIDQuery struct {
	s personqs.QueryService
}

func NewListByOrganizationIDQuery(s personqs.QueryService) *ListByOrganizationIDQuery {
	return &ListByOrganizationIDQuery{
		s: s,
	}
}

type ListByOrganizationIDInput struct {
	OrganizationID string
	Limit          int
	Offset         int
}

type ListByOrganizationIDOutput struct {
	Views []personqs.PersonShortView
}

func (q *ListByOrganizationIDQuery) Execute(ctx context.Context, in ListByOrganizationIDInput) (*ListByOrganizationIDOutput, error) {
	orgID, err := uuid.Parse(in.OrganizationID)
	if err != nil {
		return nil, err
	}

	views, err := q.s.ListByOrganizationID(ctx, orgID, in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	return &ListByOrganizationIDOutput{
		Views: views,
	}, nil
}
