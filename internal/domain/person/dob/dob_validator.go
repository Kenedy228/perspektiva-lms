package dob

import (
	"errors"
	"time"
)

var (
	ErrInvalidDob = errors.New("invalid date of birth")
	ErrInvalidAt  = errors.New("invalid date for calculating age")
)

func validateDateOfBirth(dob time.Time) error {
	if time.Now().Before(dob) {
		return ErrInvalidDob
	}

	return nil
}

func validateForAgeAt(dob, at time.Time) error {
	if at.Before(dob) {
		return ErrInvalidAt
	}

	return nil
}
