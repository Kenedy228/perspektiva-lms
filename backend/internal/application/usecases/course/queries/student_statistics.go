package queries

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type StudentStatisticsQuery struct {
	s courseports.QueryService
}

func NewStudentStatisticsQuery(s courseports.QueryService) *StudentStatisticsQuery {
	if s == nil {
		panic("student statistics query requires query service")
	}
	return &StudentStatisticsQuery{s: s}
}

type StudentStatisticsInput struct {
	ActorRole      role.Role
	AccountID      string
	OrganizationID string
	Limit          int
	Offset         int
}

type StudentStatisticsOutput struct {
	Views []courseports.StudentRatingView
}

func (q *StudentStatisticsQuery) Execute(ctx context.Context, in StudentStatisticsInput) (*StudentStatisticsOutput, error) {
	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}

	filter := courseports.StudentStatisticsFilter{
		Limit:  limit,
		Offset: offset,
	}

	switch in.ActorRole.Kind() {
	case role.TypeAdmin:
	case role.TypeStudent:
		accountID, err := parseRequiredUUID(in.AccountID, "account id")
		if err != nil {
			return nil, err
		}
		filter.AccountID = accountID
	case role.TypeOrganization:
		organizationID, err := parseRequiredUUID(in.OrganizationID, "organization id")
		if err != nil {
			return nil, err
		}
		filter.OrganizationID = organizationID
	case role.TypeCreator:
		return nil, fmt.Errorf("%w: creator cannot view student statistics", common.ErrForbidden)
	default:
		return nil, fmt.Errorf("%w: unsupported role", common.ErrForbidden)
	}

	views, err := q.s.ListStudentStatistics(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list student statistics: %w", err)
	}

	return &StudentStatisticsOutput{Views: views}, nil
}

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}
	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: %s is required", common.ErrInvalidInput, field)
	}
	return id, nil
}
