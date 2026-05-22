package course

import (
	"context"

	"github.com/google/uuid"
)

type Filter struct {
	TitleContains string
	Status        string
	AccountID     uuid.UUID
	Limit         int
	Offset        int
}

type QueryService interface {
	ListManageable(ctx context.Context, filter Filter) ([]ShortView, error)
	ListVisibleForStudent(ctx context.Context, filter Filter) ([]ShortView, error)
	GetDetailsByID(ctx context.Context, id uuid.UUID) (DetailedView, error)
	ListRatings(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]StudentRatingView, error)
	ListStudentStatistics(ctx context.Context, filter StudentStatisticsFilter) ([]StudentRatingView, error)
}

type StudentStatisticsFilter struct {
	AccountID      uuid.UUID
	OrganizationID uuid.UUID
	Limit          int
	Offset         int
}
