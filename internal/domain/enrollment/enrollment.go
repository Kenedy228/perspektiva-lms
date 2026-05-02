package enrollment

import (
	"time"

	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Enrollment struct {
	id            uuid.UUID
	courseID      uuid.UUID
	versionID     uuid.UUID
	accountID     uuid.UUID
	activatedAt   time.Time
	deactivatedAt time.Time
}

func New(
	courseID uuid.UUID,
	versionID uuid.UUID,
	studentID uuid.UUID,
	activatedAt time.Time,
	deactivatedAt time.Time,
) (*Enrollment, error) {
	activatedAt = normalizeDate(activatedAt)
	deactivatedAt = normalizeDate(deactivatedAt)

	if err := validateRequiredID("courseID", courseID); err != nil {
		return nil, err
	}

	if err := validateRequiredID("versionID", versionID); err != nil {
		return nil, err
	}

	if err := validateRequiredID("studentID", studentID); err != nil {
		return nil, err
	}

	if err := validateActivationWindow(activatedAt, deactivatedAt); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Enrollment{
		id:            id,
		courseID:      courseID,
		versionID:     versionID,
		accountID:     studentID,
		activatedAt:   activatedAt,
		deactivatedAt: deactivatedAt,
	}, nil
}

func (e *Enrollment) ID() uuid.UUID {
	return e.id
}

func (e *Enrollment) CourseID() uuid.UUID {
	return e.courseID
}

func (e *Enrollment) VersionID() uuid.UUID {
	return e.versionID
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

func (e *Enrollment) Status(at time.Time) Status {
	today := normalizeDate(at)

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

func (e *Enrollment) ChangeActivationWindow(activatedAt, deactivatedAt time.Time) error {
	activatedAt = normalizeDate(activatedAt)
	deactivatedAt = normalizeDate(deactivatedAt)

	if err := validateActivationWindow(activatedAt, deactivatedAt); err != nil {
		return err
	}

	e.activatedAt = activatedAt
	e.deactivatedAt = deactivatedAt

	return nil
}

func (e *Enrollment) Clone() *Enrollment {
	return &Enrollment{
		id:            e.id,
		courseID:      e.courseID,
		versionID:     e.versionID,
		accountID:     e.accountID,
		activatedAt:   e.activatedAt,
		deactivatedAt: e.deactivatedAt,
	}
}
