package enrollment

import (
	"context"

	enrollmentdomain "gitflic.ru/lms/backend/internal/domain/enrollment"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*enrollmentdomain.Enrollment, error)
	Save(ctx context.Context, e *enrollmentdomain.Enrollment) error
	ExistsForAccountCourseVersion(ctx context.Context, accountID, courseID, versionID uuid.UUID, excludeEnrollmentID uuid.UUID) (bool, error)
}

type VersionPolicy interface {
	CanEnrollVersion(ctx context.Context, versionID uuid.UUID) (bool, error)
}

type ProgressInitializer interface {
	EnsureProgressForEnrollment(ctx context.Context, enrollmentID uuid.UUID) error
}
