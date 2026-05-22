package queries

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	"gitflic.ru/lms/backend/internal/domain/role"
)

type GetByIDQuery struct {
	s accountports.QueryService
}

func NewGetByIDQuery(s accountports.QueryService) *GetByIDQuery {
	if s == nil {
		panic("account get by id query requires query service")
	}

	return &GetByIDQuery{s: s}
}

type GetByIDInput struct {
	ActorRole role.Role
	AccountID string
}

type GetOutput struct {
	View accountports.AccountView
}

func (q *GetByIDQuery) Execute(ctx context.Context, in GetByIDInput) (*GetOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(in.AccountID, "account id")
	if err != nil {
		return nil, err
	}

	view, err := q.s.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get account by id: %w", err)
	}

	return &GetOutput{View: view}, nil
}

type GetByPersonIDQuery struct {
	s accountports.QueryService
}

func NewGetByPersonIDQuery(s accountports.QueryService) *GetByPersonIDQuery {
	if s == nil {
		panic("account get by person id query requires query service")
	}

	return &GetByPersonIDQuery{s: s}
}

type GetByPersonIDInput struct {
	ActorRole role.Role
	PersonID  string
}

func (q *GetByPersonIDQuery) Execute(ctx context.Context, in GetByPersonIDInput) (*GetOutput, error) {
	if err := common.RequireAdmin(in.ActorRole); err != nil {
		return nil, err
	}

	personID, err := parseRequiredUUID(in.PersonID, "person id")
	if err != nil {
		return nil, err
	}

	view, err := q.s.GetByPersonID(ctx, personID)
	if err != nil {
		return nil, fmt.Errorf("get account by person id: %w", err)
	}

	return &GetOutput{View: view}, nil
}
