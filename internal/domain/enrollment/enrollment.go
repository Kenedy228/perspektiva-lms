package enrollment

import (
	"time"

	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Enrollment struct {
	id              uuid.UUID
	courseID        uuid.UUID
	courseVersionID uuid.UUID
	accountID       uuid.UUID
	activatedAt     time.Time
	deactivatedAt   time.Time
	createdAt       time.Time
	updatedAt       time.Time
}

func New(params Params) (*Enrollment, error) {
	activatedAt := normalize(params.ActivatedAt)
	deactivatedAt := normalize(params.DeactivatedAt)

	if err := validateCourseID(params.CourseID); err != nil {
		return nil, err
	}

	if err := validateCourseVersionID(params.CourseVersionID); err != nil {
		return nil, err
	}

	if err := validateAccountID(params.AccountID); err != nil {
		return nil, err
	}

	if err := validateTimeBoundaries(activatedAt, deactivatedAt); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	return &Enrollment{
		id:              id,
		courseID:        params.CourseID,
		courseVersionID: params.CourseVersionID,
		accountID:       params.AccountID,
		activatedAt:     activatedAt,
		deactivatedAt:   deactivatedAt,
		createdAt:       now,
		updatedAt:       now,
	}, nil
}

func (e *Enrollment) ID() uuid.UUID {
	return e.id
}

func (e *Enrollment) CourseID() uuid.UUID {
	return e.courseID
}

func (e *Enrollment) CourseVersionID() uuid.UUID {
	return e.courseVersionID
}

func (e *Enrollment) AccountID() uuid.UUID {
	return e.accountID
}

func (e *Enrollment) ActivatedAt() time.Time {
	return e.activatedAt
}

func (e *Enrollment) DeactivatedAt() time.Time {
	return e.deactivatedAt
}

func (e *Enrollment) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Enrollment) UpdatedAt() time.Time {
	return e.updatedAt
}

func (e *Enrollment) Status(at time.Time) Status {
	today := normalize(at)

	if today.Before(e.activatedAt) {
		return StatusInactive
	}

	if today.After(e.deactivatedAt) {
		return StatusExpired
	}

	return StatusActive
}

func (e *Enrollment) IsActive(at time.Time) bool {
	return e.Status(at) == StatusActive
}

func (e *Enrollment) ChangeActivationWindow(from, to time.Time) error {
	from = normalize(from)
	to = normalize(to)

	if err := validateTimeBoundaries(from, to); err != nil {
		return err
	}

	e.activatedAt = from
	e.deactivatedAt = to
	e.updatedAt = time.Now()

	return nil
}
