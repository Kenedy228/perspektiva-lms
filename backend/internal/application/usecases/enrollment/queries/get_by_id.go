package queries

import (
	"context"
	"fmt"

	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	"gitflic.ru/lms/backend/internal/application/usecases/enrollment/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type GetByIDQuery struct {
	service enrollmentports.QueryService
}

func NewGetByIDQuery(service enrollmentports.QueryService) *GetByIDQuery {
	if service == nil {
		panic("enrollment get by id query requires query service")
	}
	return &GetByIDQuery{service: service}
}

type GetByIDInput struct {
	ActorRole      role.Role
	ActorAccountID string
	EnrollmentID   string
}

type GetByIDOutput struct {
	View enrollmentports.EnrollmentView
}

func (q *GetByIDQuery) Execute(ctx context.Context, in GetByIDInput) (*GetByIDOutput, error) {
	if err := common.RequireViewer(in.ActorRole); err != nil {
		return nil, err
	}

	id, err := parseRequiredUUID(in.EnrollmentID, "enrollment id")
	if err != nil {
		return nil, err
	}

	view, err := q.service.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get enrollment by id: %w", err)
	}

	if in.ActorRole.Kind() == role.TypeStudent {
		actorAccountID, err := parseRequiredUUID(in.ActorAccountID, "actor account id")
		if err != nil {
			return nil, err
		}
		if view.AccountID != actorAccountID.String() {
			return nil, fmt.Errorf("%w: студент может просматривать только свои зачисления", common.ErrForbidden)
		}
	}

	return &GetByIDOutput{View: view}, nil
}

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: поле '%s' обязательно", common.ErrInvalidInput, field)
	}
	return id, nil
}
