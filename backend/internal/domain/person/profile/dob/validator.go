package dob

import (
	"fmt"
	"time"

	calculator "github.com/Kenedy228/age-calculator"
)

// validateAdultDateOfBirth вычисляет возраст через AgeAt и проверяет, что он находится
// в допустимых границах взрослого возраста.
func validateAdultDateOfBirth(date, asOf time.Time) error {
	age, err := calculator.AgeAt(date, asOf)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	if err := validateMinAgeBoundary(age); err != nil {
		return err
	}

	if err := validateMaxAgeBoundary(age); err != nil {
		return err
	}

	return nil
}

// validateMinAgeBoundary проверяет, что возраст не ниже минимального порога взрослого возраста.
func validateMinAgeBoundary(age int) error {
	if age < MinAdultAge {
		return fmt.Errorf("%w: минимальный возраст человека должен быть не менее %d лет, получено %d", ErrInvalid, MinAdultAge, age)
	}

	return nil
}

// validateMaxAgeBoundary проверяет, что возраст не превышает максимальный допустимый порог.
func validateMaxAgeBoundary(age int) error {
	if age > MaxAdultAge {
		return fmt.Errorf("%w: максимальный возраст человека должен быть не более %d лет, получено %d", ErrInvalid, MaxAdultAge, age)
	}

	return nil
}
