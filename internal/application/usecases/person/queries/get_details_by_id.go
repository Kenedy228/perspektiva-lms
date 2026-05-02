package queries

import (
	"context"

	personqs "gitflic.ru/lms/internal/application/ports/person"
	"github.com/google/uuid"
)

type GetDetailsByIdQuery struct {
	s personqs.QueryService
}

func NewGetDetailsByIdQuery(s personqs.QueryService) *GetDetailsByIdQuery {
	return &GetDetailsByIdQuery{
		s: s,
	}
}

type GetDetailsByIDInput struct {
	ID string
}

type GetDetailsByIDOutput struct {
	View personqs.PersonDetailedView
}

func (qs *GetDetailsByIdQuery) Execute(ctx context.Context, in GetDetailsByIDInput) (*GetDetailsByIDOutput, error) {
	id, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, err
	}

	view, err := qs.s.GetDetailsByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &GetDetailsByIDOutput{
		View: view,
	}, nil
}
