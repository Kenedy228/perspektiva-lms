package jobtitle

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func validateJobTitle(jobTitle string) error {
	if err := validateRequiredJobTitle(jobTitle); err != nil {
		return err
	}

	if err := validateCharsLimitJobTitle(jobTitle); err != nil {
		return err
	}

	return nil
}

func validateRequiredJobTitle(jobTitle string) error {
	if strings.TrimSpace(jobTitle) == "" {
		return fmt.Errorf("%w, детали: должность должна содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateCharsLimitJobTitle(jobTitle string) error {
	if utf8.RuneCountInString(jobTitle) > JobTitleCharsLimit {
		return fmt.Errorf("%w, детали: должность должна содержать не более %d символов", ErrInvalid, JobTitleCharsLimit)
	}

	return nil
}
