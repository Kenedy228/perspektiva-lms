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
	_ enrollmentports.OrganizationScope   = (*EnrollmentRepository)(nil)
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
	var courseID, accountID uuid.UUID
	var activatedAt, deactivatedAt time.Time
	err := r.db.QueryRowContext(ctx, `
		SELECT course_id, account_id, enrolled_at, completed_at
		FROM enrollments
		WHERE id = $1`, id).Scan(&courseID, &accountID, &activatedAt, &deactivatedAt)
	if err != nil {
		return nil, err
	}
	return enrollmentdomain.Restore(id, courseID, accountID, activatedAt, deactivatedAt)
}

func (r *EnrollmentRepository) Save(ctx context.Context, e *enrollmentdomain.Enrollment) error {
	nowStatus := e.Status(time.Now()).String()
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO enrollments (
			id, account_id, course_id, status, enrolled_at, completed_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, now())
		ON CONFLICT (id) DO UPDATE SET
			account_id = EXCLUDED.account_id,
			course_id = EXCLUDED.course_id,
			status = EXCLUDED.status,
			enrolled_at = EXCLUDED.enrolled_at,
			completed_at = EXCLUDED.completed_at,
			updated_at = now()`,
		e.ID(), e.AccountID(), e.CourseID(), nowStatus, e.ActivatedAt(), e.DeactivatedAt())
	if err != nil {
		return fmt.Errorf("save enrollment: %w", err)
	}
	return nil
}

func (r *EnrollmentRepository) ExistsForAccountCourse(ctx context.Context, accountID, courseID uuid.UUID, excludeEnrollmentID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM enrollments
			WHERE account_id = $1 AND course_id = $2
				AND ($3::uuid IS NULL OR id <> $3)
		)`, accountID, courseID, nullUUID(excludeEnrollmentID)).Scan(&exists)
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

func (r *EnrollmentRepository) EnrollmentBelongsToPersonOrganization(ctx context.Context, enrollmentID, personID uuid.UUID) (bool, error) {
	var ok bool
	err := r.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM enrollments e
			JOIN accounts a ON a.id = e.account_id
			JOIN persons pe ON pe.id = a.person_id
			JOIN persons pm ON pm.id = $2
			WHERE e.id = $1
				AND pe.deleted_at IS NULL
				AND pm.deleted_at IS NULL
				AND pe.organization_id IS NOT NULL
				AND pm.organization_id IS NOT NULL
				AND pe.organization_id = pm.organization_id
		)`, enrollmentID, personID).Scan(&ok)
	if err != nil {
		return false, fmt.Errorf("check enrollment organization scope: %w", err)
	}
	return ok, nil
}

func (r *EnrollmentRepository) PersonOrganizationID(ctx context.Context, personID uuid.UUID) (uuid.UUID, error) {
	var organizationID uuid.NullUUID
	err := r.db.QueryRowContext(ctx, `
		SELECT organization_id
		FROM persons
		WHERE id = $1 AND deleted_at IS NULL`, personID).Scan(&organizationID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("lookup person organization: %w", err)
	}
	if !organizationID.Valid {
		return uuid.Nil, nil
	}
	return organizationID.UUID, nil
}
