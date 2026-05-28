package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	enrollmentdomain "gitflic.ru/lms/backend/internal/domain/enrollment"
	"github.com/google/uuid"
)

var (
	_ enrollmentports.Repository          = (*EnrollmentRepository)(nil)
	_ enrollmentports.ProgressInitializer = (*EnrollmentRepository)(nil)
)

// EnrollmentRepository persists enrollment aggregates and initializes progress rows.
type EnrollmentRepository struct {
	db *sql.DB
}

// NewEnrollmentRepository creates a PostgreSQL enrollment adapter.
func NewEnrollmentRepository(db *sql.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

func (r *EnrollmentRepository) FindByID(ctx context.Context, id uuid.UUID) (*enrollmentdomain.Enrollment, error) {
	var courseID, versionID, accountID uuid.UUID
	var activatedAt, deactivatedAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT course_id, version_id, account_id, enrolled_at, completed_at
		FROM enrollments
		WHERE id = $1`, id).Scan(&courseID, &versionID, &accountID, &activatedAt, &deactivatedAt)
	if err != nil {
		return nil, err
	}
	return enrollmentdomain.Restore(id, courseID, versionID, accountID, activatedAt, deactivatedAt)
}

func (r *EnrollmentRepository) Save(ctx context.Context, e *enrollmentdomain.Enrollment) error {
	nowStatus := e.Status(time.Now()).String()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO enrollments (
			id, account_id, course_id, version_id, status, enrolled_at, completed_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, now())
		ON CONFLICT (id) DO UPDATE SET
			account_id = EXCLUDED.account_id,
			course_id = EXCLUDED.course_id,
			version_id = EXCLUDED.version_id,
			status = EXCLUDED.status,
			enrolled_at = EXCLUDED.enrolled_at,
			completed_at = EXCLUDED.completed_at,
			updated_at = now()`,
		e.ID(), e.AccountID(), e.CourseID(), e.VersionID(), nowStatus, e.ActivatedAt(), e.DeactivatedAt())
	if err != nil {
		return fmt.Errorf("save enrollment: %w", err)
	}
	return nil
}

func (r *EnrollmentRepository) ExistsForAccountCourseVersion(ctx context.Context, accountID, courseID, versionID uuid.UUID, excludeEnrollmentID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM enrollments
			WHERE account_id = $1 AND course_id = $2 AND version_id = $3
				AND ($4::uuid IS NULL OR id <> $4)
		)`, accountID, courseID, versionID, nullUUID(excludeEnrollmentID)).Scan(&exists)
	return exists, err
}

func (r *EnrollmentRepository) EnsureProgressForEnrollment(ctx context.Context, enrollmentID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO course_progress (enrollment_id)
		VALUES ($1)
		ON CONFLICT (enrollment_id) DO NOTHING`, enrollmentID)
	if err != nil {
		return fmt.Errorf("ensure progress for enrollment: %w", err)
	}
	return nil
}
