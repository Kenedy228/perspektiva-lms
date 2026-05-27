package score

import "fmt"

func validateValue(value float64) error {
	if err := validateNegativeValue(value); err != nil {
		return err
	}

	if err := validateMaxValue(value); err != nil {
		return err
	}

	return nil
}

func validateNegativeValue(value float64) error {
	if value < 0 {
		return fmt.Errorf("%w: значение меньше нуля", ErrInvalid)
	}

	return nil
}

func validateMaxValue(value float64) error {
	if value > 1 {
		return fmt.Errorf("%w: значение больше единицы", ErrInvalid)
	}

	return nil
}
