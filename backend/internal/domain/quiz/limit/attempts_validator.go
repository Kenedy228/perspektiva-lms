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
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if count > maxAttemptsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, maxAttemptsCount)
	}

	return nil
}
