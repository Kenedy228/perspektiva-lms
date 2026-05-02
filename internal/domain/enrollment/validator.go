package enrollment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func validateRequiredID(fieldName string, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("%w, детали: %s должен быть валидным идентификатором", ErrInvalid, fieldName)
	}
	return nil
}

func validateActivationWindow(activatedAt, deactivatedAt time.Time) error {
	today := normalizeDate(time.Now())

	if activatedAt.Before(today) {
		return fmt.Errorf("%w, детали: дата активации не может быть раньше сегодняшнего дня", ErrInvalid)
	}

	if deactivatedAt.Before(activatedAt) {
		return fmt.Errorf("%w, детали: дата деактивации не может быть раньше даты активации", ErrInvalid)
	}

	return nil
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
