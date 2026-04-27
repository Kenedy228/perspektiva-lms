package dob

import (
	"time"
)

func validateAdultDateOfBirth(date, asOf time.Time) error {
	if err := validateNotInFuture(date, asOf); err != nil {
		return err
	}

	if err := validateAdultAge(date, asOf); err != nil {
		return err
	}

	return nil
}
