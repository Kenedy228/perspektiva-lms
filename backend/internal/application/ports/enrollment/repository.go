package enrollment

import (
	"context"

	enrollmentdomain "gitflic.ru/lms/backend/internal/domain/enrollment"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*enrollmentdomain.Enrollment, error)
	Save(ctx context.Context, e *enrollmentdomain.Enrollment) error
	ExistsForAccountCourse(ctx context.Context, accountID, courseID uuid.UUID, excludeEnrollmentID uuid.UUID) (bool, error)
}

type ProgressInitializer interface {
	EnsureProgressForEnrollment(ctx context.Context, enrollmentID uuid.UUID) error
}

type OrganizationScope interface {
	EnrollmentBelongsToPersonOrganization(ctx context.Context, enrollmentID, personID uuid.UUID) (bool, error)
	PersonOrganizationID(ctx context.Context, personID uuid.UUID) (uuid.UUID, error)
}
