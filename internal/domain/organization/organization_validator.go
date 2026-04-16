package organization

import (
	"errors"
	"regexp"
	"strings"
)

var (
	individualRegexp = regexp.MustCompile(`^\d{12}$`)
	companyRegexp    = regexp.MustCompile(`^\d{10}$`)
)

var (
	ErrEmptyName  = errors.New("empty name")
	ErrEmptyInn   = errors.New("empty inn")
	ErrInvalidInn = errors.New("invalid inn")
)

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}

	return nil
}

func validateInn(inn string) error {
	if strings.TrimSpace(inn) == "" {
		return ErrEmptyInn
	}

	if !individualRegexp.MatchString(inn) && !companyRegexp.MatchString(inn) {
		return ErrInvalidInn
	}

	return nil
}
