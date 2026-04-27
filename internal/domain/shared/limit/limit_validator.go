package limit

import "fmt"

func validateValue(value int) error {
	if value < 0 {
		return fmt.Errorf("%w, детали: значение лимита не может быть отрицательным", ErrInvalidLimit)
	}

	if value > maxValueInSecond {
		return fmt.Errorf("%w, детали: значение лимита не должно превышать 1 года", ErrInvalidLimit)
	}

	return nil
}
