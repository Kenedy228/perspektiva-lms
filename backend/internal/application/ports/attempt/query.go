package attempt

import (
	"context"

	"github.com/google/uuid"
)

type QueryService interface {
	ListByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID, limit, offset int) ([]AttemptSummaryView, error)
}
