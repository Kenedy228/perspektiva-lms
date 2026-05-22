package dob

import (
	"fmt"
	"time"

	calculator "github.com/Kenedy228/age-calculator"
)

func validateAdultDateOfBirth(date, asOf time.Time) error {
	age, err := calculator.AgeAt(date, asOf)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	if age < MinAdultAge {
		return fmt.Errorf("%w: минимальный возраст человека должен быть не менее %d лет", ErrInvalid, MinAdultAge)
	}

	if age > MaxAdultAge {
		return fmt.Errorf("%w: максимальный возраст человека должен быть не более %d лет", ErrInvalid, MaxAdultAge)
	}

	return nil
}
