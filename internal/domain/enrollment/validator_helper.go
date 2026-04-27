package enrollment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func validateRequiredID(field string, value uuid.UUID) error {
	if value == uuid.Nil {
		return fmt.Errorf("%w, детали: %s должен быть существующим объектом", ErrInvalid, field)
	}

	return nil
}

func validateDateNotBefore(field string, date, at time.Time) error {
	if date.Before(at) {
		return fmt.Errorf("%w, детали: %s не может быть раньше %s", ErrInvalid, field, at.Format("2006-01-02"))
	}
	return nil
}
