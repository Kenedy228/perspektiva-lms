package organization

import (
	"fmt"
	"strings"
)

func validateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("%w, детали: наименование организации должно содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateInn(inn string) error {
	if !individualRegexp.MatchString(inn) && !companyRegexp.MatchString(inn) {
		return fmt.Errorf("%w, детали: предоставленный формат ИНН не является корректным", ErrInvalid)
	}

	return nil
}
