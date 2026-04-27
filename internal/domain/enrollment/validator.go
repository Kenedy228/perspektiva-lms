package enrollment

import (
	"time"

	"github.com/google/uuid"
)

func validateCourseID(id uuid.UUID) error {
	return validateRequiredID(courseIDField, id)
}

func validateCourseVersionID(id uuid.UUID) error {
	return validateRequiredID(courseVersionIDField, id)
}

func validateAccountID(id uuid.UUID) error {
	return validateRequiredID(accountIDField, id)
}

func validateTimeBoundaries(from, to time.Time) error {
	today := normalize(time.Now())

	if err := validateDateNotBefore("дата активации", from, today); err != nil {
		return err
	}

	if err := validateDateNotBefore("дата деактивации", to, from); err != nil {
		return err
	}

	return nil
}
