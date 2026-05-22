package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/person"
	"gitflic.ru/lms/backend/internal/application/usecases/person/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type GetDetailsByIdQuery struct {
	s person.QueryService
}

func NewGetDetailsByIdQuery(s person.QueryService) *GetDetailsByIdQuery {
	if s == nil {
		panic("person get details query requires query service")
	}

	return &GetDetailsByIdQuery{
		s: s,
	}
}

type GetDetailsByIDInput struct {
	ActorRole role.Role
	ID        string
}

type GetDetailsByIDOutput struct {
	View person.PersonDetailedView
}

func (q *GetDetailsByIdQuery) Execute(ctx context.Context, in GetDetailsByIDInput) (*GetDetailsByIDOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(in.ID, "person id")
	if err != nil {
		return nil, err
	}

	view, err := q.s.GetDetailsByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get person details: %w", err)
	}

	return &GetDetailsByIDOutput{
		View: view,
	}, nil
}
