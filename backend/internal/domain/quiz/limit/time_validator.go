package limit

import "fmt"

func validateSeconds(seconds int) error {
	if err := validateBoundaries(seconds); err != nil {
		return err
	}

	return nil
}

func validateBoundaries(seconds int) error {
	if seconds < 0 {
		return fmt.Errorf("%w: invalid value", ErrInvalid)
	}

	if seconds > maxSecondsCount {
		return fmt.Errorf("%w: invalid value (%d)", ErrInvalid, maxSecondsCount)
	}

	return nil
}
