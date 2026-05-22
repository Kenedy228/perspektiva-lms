//go:build legacy
// +build legacy

package enrollment_test

import (
	"testing"
	"time"

	enrollment2 "gitflic.ru/lms/backend/internal/domain/enrollment"
	"gitflic.ru/lms/backend/internal/domain/enrollment"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type enrollmentBuilder struct {
	courseID        uuid.UUID
	courseVersionID uuid.UUID
	accountID       uuid.UUID
	activatedAt     time.Time
	deactivatedAt   time.Time
}

func newEnrollmentBuilder() *enrollmentBuilder {
	today := time.Now()
	start := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 0, 30)

	return &enrollmentBuilder{
		courseID:        uuid.New(),
		courseVersionID: uuid.New(),
		accountID:       uuid.New(),
		activatedAt:     start,
		deactivatedAt:   end,
	}
}

func (b *enrollmentBuilder) withCourseID(id uuid.UUID) *enrollmentBuilder {
	b.courseID = id
	return b
}

func (b *enrollmentBuilder) withCourseVersionID(id uuid.UUID) *enrollmentBuilder {
	b.courseVersionID = id
	return b
}

func (b *enrollmentBuilder) withAccountID(id uuid.UUID) *enrollmentBuilder {
	b.accountID = id
	return b
}

func (b *enrollmentBuilder) withWindow(from, to time.Time) *enrollmentBuilder {
	b.activatedAt = from
	b.deactivatedAt = to
	return b
}

func (b *enrollmentBuilder) build(t *testing.T, wantErr error) *enrollment2.Enrollment {
	t.Helper()

	params := enrollment.Params{
		CourseID:        b.courseID,
		CourseVersionID: b.courseVersionID,
		AccountID:       b.accountID,
		ActivatedAt:     b.activatedAt,
		DeactivatedAt:   b.deactivatedAt,
	}

	e, err := enrollment2.New(params)
	assert.ErrorIs(t, err, wantErr)

	return e
}
