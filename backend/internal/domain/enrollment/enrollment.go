package enrollment

import (
	"fmt"
	"time"

	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Enrollment struct {
	id            uuid.UUID
	courseID      uuid.UUID
	accountID     uuid.UUID
	activatedAt   time.Time
	deactivatedAt time.Time
}

func New(
	courseID uuid.UUID,
	accountID uuid.UUID,
	activatedAt time.Time,
	deactivatedAt time.Time,
) (*Enrollment, error) {
	return NewAt(courseID, accountID, activatedAt, deactivatedAt, time.Now())
}

func NewAt(
	courseID uuid.UUID,
	accountID uuid.UUID,
	activatedAt time.Time,
	deactivatedAt time.Time,
	now time.Time,
) (*Enrollment, error) {
	activatedAt = normalizeDate(activatedAt)
	deactivatedAt = normalizeDate(deactivatedAt)
	now = normalizeDate(now)

	if err := validateRequiredID("courseID", courseID); err != nil {
		return nil, err
	}

	if err := validateRequiredID("accountID", accountID); err != nil {
		return nil, err
	}

	if err := validateActivationWindowAt(activatedAt, deactivatedAt, now); err != nil {
		return nil, err
	}

	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Enrollment{
		id:            id,
		courseID:      courseID,
		accountID:     accountID,
		activatedAt:   activatedAt,
		deactivatedAt: deactivatedAt,
	}, nil
}

func Restore(
	id uuid.UUID,
	courseID uuid.UUID,
	accountID uuid.UUID,
	activatedAt time.Time,
	deactivatedAt time.Time,
) (*Enrollment, error) {
	activatedAt = normalizeDate(activatedAt)
	deactivatedAt = normalizeDate(deactivatedAt)

	if err := validateRequiredID("id", id); err != nil {
		return nil, err
	}
	if err := validateRequiredID("courseID", courseID); err != nil {
		return nil, err
	}
	if err := validateRequiredID("accountID", accountID); err != nil {
		return nil, err
	}
	if err := validateActivationWindowOrder(activatedAt, deactivatedAt); err != nil {
		return nil, err
	}

	return &Enrollment{
		id:            id,
		courseID:      courseID,
		accountID:     accountID,
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
	return e.ChangeActivationWindowAt(activatedAt, deactivatedAt, time.Now())
}

func (e *Enrollment) ChangeActivationWindowAt(activatedAt, deactivatedAt time.Time, now time.Time) error {
	activatedAt = normalizeDate(activatedAt)
	deactivatedAt = normalizeDate(deactivatedAt)
	now = normalizeDate(now)

	if err := validateActivationWindowAt(activatedAt, deactivatedAt, now); err != nil {
		return err
	}

	e.activatedAt = activatedAt
	e.deactivatedAt = deactivatedAt

	return nil
}

func (e *Enrollment) Deactivate(at time.Time) error {
	at = normalizeDate(at)
	if at.Before(e.activatedAt) {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}
	e.deactivatedAt = at
	return nil
}

func (e *Enrollment) Clone() *Enrollment {
	return &Enrollment{
		id:            e.id,
		courseID:      e.courseID,
		accountID:     e.accountID,
		activatedAt:   e.activatedAt,
		deactivatedAt: e.deactivatedAt,
	}
}
