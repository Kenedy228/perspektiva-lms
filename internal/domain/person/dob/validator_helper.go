package dob

import (
	"fmt"
	"time"
)

func validateNotInFuture(date, asOf time.Time) error {
	if asOf.Before(date) {
		return fmt.Errorf("%w, детали: дата рождения не может быть в будущем", ErrInvalid)
	}

	return nil
}

func validateAdultAge(date, asOf time.Time) error {
	if ageAt(date, asOf) < minAdultAgeYears {
		return fmt.Errorf("%w, детали: добавляемый студент младше %d лет", ErrInvalid, minAdultAgeYears)
	}

	return nil
}
