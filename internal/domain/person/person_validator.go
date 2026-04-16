package person

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrEmptyFirstName    = errors.New("empty firstName")
	ErrEmptyLastName     = errors.New("empty lastName")
	ErrEmptyMiddleName   = errors.New("empty middleName")
	ErrInvalidFirstName  = errors.New("invalid firstName")
	ErrInvalidLastName   = errors.New("invalid lastName")
	ErrInvalidMiddleName = errors.New("invalid middleName")
	ErrEmptyJobTitle     = errors.New("empty job title")
	ErrInvalidSnils      = errors.New("invalid snils")
	ErrEmptyEducation    = errors.New("empty education")
)

var (
	nameRegex  = regexp.MustCompile(`^[А-Яа-яЁё\s\-']+$`)
	snilsRegex = regexp.MustCompile(`^\d{3}-\d{3}-\d{3} \d{2}$`)
)

func validateFirstName(firstName string) error {
	empty, invalid := validateNamePart(firstName)
	if empty {
		return ErrEmptyFirstName
	}
	if invalid {
		return ErrInvalidFirstName
	}

	return nil
}

func validateLastName(lastName string) error {
	empty, invalid := validateNamePart(lastName)
	if empty {
		return ErrEmptyLastName
	}
	if invalid {
		return ErrInvalidLastName
	}

	return nil
}

func validateMiddleName(middleName string) error {
	empty, invalid := validateNamePart(middleName)
	if empty {
		return nil
	}

	if invalid {
		return ErrInvalidMiddleName
	}

	return nil
}

func validateJobTitle(jobTitle string) error {
	if strings.TrimSpace(jobTitle) == "" {
		return ErrEmptyJobTitle
	}

	return nil
}

func validateSnils(snils string) error {
	if !snilsRegex.MatchString(snils) {
		return ErrInvalidSnils
	}

	return nil
}

func validateEducation(education string) error {
	if strings.TrimSpace(education) == "" {
		return ErrEmptyEducation
	}
	return nil
}

func validateNamePart(val string) (empty bool, invalid bool) {
	if strings.TrimSpace(val) == "" {
		empty = true
		return
	}

	if !nameRegex.MatchString(val) {
		invalid = true
		return
	}
	return
}
