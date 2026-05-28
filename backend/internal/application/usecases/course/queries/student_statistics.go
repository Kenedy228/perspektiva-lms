package queries

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type StudentStatisticsQuery struct {
	s        courseports.QueryService
	orgScope enrollmentports.OrganizationScope
}

func NewStudentStatisticsQuery(s courseports.QueryService, orgScope enrollmentports.OrganizationScope) *StudentStatisticsQuery {
	if s == nil {
		panic("student statistics query requires query service")
	}
	return &StudentStatisticsQuery{s: s, orgScope: orgScope}
}

type StudentStatisticsInput struct {
	ActorRole      role.Role
	ActorPersonID  string
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
		personID, err := parseRequiredUUID(in.ActorPersonID, "actor person id")
		if err != nil {
			return nil, err
		}
		actorOrgID, err := q.orgScope.PersonOrganizationID(ctx, personID)
		if err != nil {
			return nil, fmt.Errorf("lookup organization for actor: %w", err)
		}
		if actorOrgID == uuid.Nil {
			return nil, fmt.Errorf("%w: пользователь не привязан ни к одной организации", common.ErrForbidden)
		}
		if in.OrganizationID != "" {
			requestedOrgID, err := parseRequiredUUID(in.OrganizationID, "organization id")
			if err != nil {
				return nil, err
			}
			if requestedOrgID != actorOrgID {
				return nil, fmt.Errorf("%w: организация может просматривать статистику только своих студентов", common.ErrForbidden)
			}
		}
		filter.OrganizationID = actorOrgID
	case role.TypeCreator:
		return nil, fmt.Errorf("%w: создатель не имеет доступа к статистике студентов", common.ErrForbidden)
	default:
		return nil, fmt.Errorf("%w: роль не поддерживается", common.ErrForbidden)
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
		return uuid.Nil, fmt.Errorf("%w: поле '%s' обязательно", common.ErrInvalidInput, field)
	}
	return id, nil
}
