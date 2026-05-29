package queries

import (
	"context"
	"fmt"

	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	"gitflic.ru/lms/backend/internal/application/usecases/enrollment/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type ListQuery struct {
	service enrollmentports.QueryService
}

func NewListQuery(service enrollmentports.QueryService) *ListQuery {
	if service == nil {
		panic("enrollment list query requires query service")
	}
	return &ListQuery{service: service}
}

type ListInput struct {
	ActorRole      role.Role
	ActorAccountID string
	AccountID      string
	CourseID       string
	Limit          int
	Offset         int
}

type ListOutput struct {
	Views []enrollmentports.EnrollmentView
}

func (q *ListQuery) Execute(ctx context.Context, in ListInput) (*ListOutput, error) {
	if err := common.RequireViewer(in.ActorRole); err != nil {
		return nil, err
	}

	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	filter := enrollmentports.ListFilter{Limit: limit, Offset: offset}

	if in.ActorRole.Kind() == role.TypeStudent {
		actorAccountID, err := parseRequiredUUID(in.ActorAccountID, "actor account id")
		if err != nil {
			return nil, err
		}
		filter.AccountID = actorAccountID
	} else {
		if in.AccountID != "" {
			accountID, err := uuid.Parse(in.AccountID)
			if err == nil && accountID != uuid.Nil {
				filter.AccountID = accountID
			}
		}
	}

	if in.CourseID != "" {
		courseID, err := uuid.Parse(in.CourseID)
		if err == nil && courseID != uuid.Nil {
			filter.CourseID = courseID
		}
	}

	views, err := q.service.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list enrollments: %w", err)
	}

	return &ListOutput{Views: views}, nil
}
