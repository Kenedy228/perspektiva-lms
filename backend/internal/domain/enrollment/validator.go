package enrollment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func validateRequiredID(fieldName string, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w: поле %s обязательно", ErrInvalid, fieldName)
	}
	return nil
}

func validateActivationWindow(activatedAt, deactivatedAt time.Time) error {
	return validateActivationWindowAt(activatedAt, deactivatedAt, normalizeDate(time.Now()))
}

func validateActivationWindowAt(activatedAt, deactivatedAt, now time.Time) error {
	if activatedAt.Before(now) {
		return fmt.Errorf("%w: дата активации не может быть в прошлом", ErrInvalid)
	}

	return validateActivationWindowOrder(activatedAt, deactivatedAt)
}

func validateActivationWindowOrder(activatedAt, deactivatedAt time.Time) error {
	if deactivatedAt.Before(activatedAt) {
		return fmt.Errorf("%w: дата деактивации не может быть раньше даты активации", ErrInvalid)
	}

	return nil
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
