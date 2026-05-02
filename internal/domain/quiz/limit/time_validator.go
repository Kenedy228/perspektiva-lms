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
		return fmt.Errorf("длительность (в секундах) лимита по времени не может быть отрицательной")
	}

	if seconds > maxSecondsCount {
		return fmt.Errorf("длительность (в секундах) лимита по времени не должна превышать %d секунд", maxSecondsCount)
	}

	return nil
}
