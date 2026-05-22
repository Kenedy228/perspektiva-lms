package queries

import (
	"context"
	"fmt"

	courseports "gitflic.ru/lms/backend/internal/application/ports/course"
	"gitflic.ru/lms/backend/internal/application/usecases/course/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

type RatingsQuery struct {
	s courseports.QueryService
}

func NewRatingsQuery(s courseports.QueryService) *RatingsQuery {
	if s == nil {
		panic("course ratings query requires query service")
	}
	return &RatingsQuery{s: s}
}

type RatingsInput struct {
	ActorRole role.Role
	CourseID  string
	Limit     int
	Offset    int
}

type RatingsOutput struct {
	Views []courseports.StudentRatingView
}

func (q *RatingsQuery) Execute(ctx context.Context, in RatingsInput) (*RatingsOutput, error) {
	if err := common.RequireManager(in.ActorRole); err != nil {
		return nil, err
	}
	limit, offset, err := common.NormalizePagination(in.Limit, in.Offset)
	if err != nil {
		return nil, err
	}
	courseID, err := uuid.Parse(in.CourseID)
	if err != nil {
		return nil, fmt.Errorf("parse course id: %w", err)
	}
	if courseID == uuid.Nil {
		return nil, fmt.Errorf("%w: course id is required", common.ErrInvalidInput)
	}
	views, err := q.s.ListRatings(ctx, courseID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list course ratings: %w", err)
	}
	return &RatingsOutput{Views: views}, nil
}
