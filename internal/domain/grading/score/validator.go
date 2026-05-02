package score

import "fmt"

func validateValue(value float64) error {
	if err := validateNegativeValue(value); err != nil {
		return err
	}

	return nil
}

func validateNegativeValue(value float64) error {
	if value < 0 {
		return fmt.Errorf("%w, детали: счет не может быть отрицательным", ErrInvalid)
	}

	return nil
}
