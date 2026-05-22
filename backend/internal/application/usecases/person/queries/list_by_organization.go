package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type ListByOrganizationIDQuery struct {
	s person.QueryService
}

func NewListByOrganizationIDQuery(s person.QueryService) *ListByOrganizationIDQuery {
	if s == nil {
		panic("person list by organization query requires query service")
	}

	return &ListByOrganizationIDQuery{
		s: s,
	}
}

type ListByOrganizationIDInput struct {
	ActorRole      role.Role
	OrganizationID string
	Limit          int
	Offset         int
}

type ListByOrganizationIDOutput struct {
	Views []person.PersonShortView
}

func (q *ListByOrganizationIDQuery) Execute(ctx context.Context, in ListByOrganizationIDInput) (*ListByOrganizationIDOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	orgID, err := parseRequiredUUID(in.OrganizationID, "organization id")
	if err != nil {
		return nil, err
	}

	views, err := q.s.ListByOrganizationID(ctx, orgID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list persons by organization: %w", err)
	}

	return &ListByOrganizationIDOutput{
		Views: views,
	}, nil
}
