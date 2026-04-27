package profile

import "fmt"

func validateJobTitle(jobTitle string) error {
	if jobTitle == "" {
		return fmt.Errorf("%w, детали: должность не может быть пустой", ErrInvalid)
	}

	if err := validateMaxLength("должность", jobTitle, jobTitleCharsLimit); err != nil {
		return err
	}

	return nil
}

func validateEducation(education string) error {
	if education == "" {
		return fmt.Errorf("%w, детали: %q не может быть пустым", ErrInvalid, "образование")
	}

	if err := validateMaxLength("образование", education, educationCharsLimit); err != nil {
		return err
	}

	return nil
}
