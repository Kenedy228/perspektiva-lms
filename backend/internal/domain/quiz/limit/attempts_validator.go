package limit

import "fmt"

func validateAttemptsCount(count int) error {
	if err := validateAttemptsCountBoundaries(count); err != nil {
		return err
	}

	return nil
}

func validateAttemptsCountBoundaries(count int) error {
	if count < 0 {
		return fmt.Errorf("%w: количество попыток не может быть отрицательным", ErrInvalid)
	}

	if count > maxAttemptsCount {
		return fmt.Errorf("%w: количество попыток не должно превышать %d", ErrInvalid, maxAttemptsCount)
	}

	return nil
}
