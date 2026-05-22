package course

import (
	"context"
	"time"

	coursedomain "gitflic.ru/lms/backend/internal/domain/course"
	"gitflic.ru/lms/backend/internal/domain/course/block"
	"gitflic.ru/lms/backend/internal/domain/course/element"
	"gitflic.ru/lms/backend/internal/domain/course/progress"
	"gitflic.ru/lms/backend/internal/domain/course/version"
	"github.com/google/uuid"
)

type CourseRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*coursedomain.Course, error)
	Save(ctx context.Context, c *coursedomain.Course) error
}

type VersionRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*version.Version, error)
	Save(ctx context.Context, v *version.Version) error
}

type BlockRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*block.Block, error)
	Save(ctx context.Context, b *block.Block) error
}

type ElementRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*element.Element, error)
	Save(ctx context.Context, e *element.Element) error
}

type ProgressRepository interface {
	FindByEnrollmentID(ctx context.Context, enrollmentID uuid.UUID) (*progress.Progress, error)
	Save(ctx context.Context, p *progress.Progress) error
}

type EnrollmentAccess interface {
	CanViewCourse(ctx context.Context, accountID, courseID uuid.UUID, at time.Time) (bool, error)
	CanEnrollVersion(ctx context.Context, versionID uuid.UUID) (bool, error)
}
