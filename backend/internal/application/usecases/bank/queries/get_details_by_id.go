package queries

import (
	"context"
	"fmt"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type GetDetailsByIDQuery struct {
	s bankports.QueryService
}

func NewGetDetailsByIDQuery(s bankports.QueryService) *GetDetailsByIDQuery {
	if s == nil {
		panic("bank get details query requires query service")
	}
	return &GetDetailsByIDQuery{s: s}
}

type GetDetailsByIDInput struct {
	ActorRole role.Role
	BankID    string
}

func (q *GetDetailsByIDQuery) Execute(ctx context.Context, in GetDetailsByIDInput) (*bankports.DetailedView, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}

	bankID, err := uuid.Parse(in.BankID)
	if err != nil {
		return nil, fmt.Errorf("parse bank id: %w", err)
	}
	if bankID == uuid.Nil {
		return nil, fmt.Errorf("%w: bank id is required", common.ErrInvalidInput)
	}

	view, err := q.s.GetDetailsByID(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("get bank details: %w", err)
	}

	return &view, nil
}
