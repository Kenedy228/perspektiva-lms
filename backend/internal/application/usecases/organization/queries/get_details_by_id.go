package queries

import (
	"context"
	"fmt"

	"gitflic.ru/lms/backend/internal/application/ports/organization"
	"gitflic.ru/lms/backend/internal/application/usecases/organization/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type GetDetailsByIDQuery struct {
	s organization.QueryService
}

func NewGetDetailsByIDQuery(s organization.QueryService) *GetDetailsByIDQuery {
	if s == nil {
		panic("organization get details query requires query service")
	}

	return &GetDetailsByIDQuery{
		s: s,
	}
}

type GetDetailsByIDInput struct {
	ActorRole role.Role
	ID        string
}

type GetDetailsByIDOutput struct {
	View organization.OrganizationDetailedView
}

func (q *GetDetailsByIDQuery) Execute(ctx context.Context, in GetDetailsByIDInput) (*GetDetailsByIDOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := uuid.Parse(in.ID)
	if err != nil {
		return nil, fmt.Errorf("parse organization id: %w", err)
	}

	view, err := q.s.GetDetailsByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get organization details: %w", err)
	}

	return &GetDetailsByIDOutput{
		View: view,
	}, nil
}
