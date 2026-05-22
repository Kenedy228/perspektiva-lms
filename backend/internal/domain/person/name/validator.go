package name

import (
	"fmt"

	"github.com/Kenedy228/fnsfio"
)

func validateFirstName(firstName string) error {
	if err := fnsfio.ValidateFirstName(firstName); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	return nil
}

func validateLastName(lastName string) error {
	if err := fnsfio.ValidateLastName(lastName); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	return nil
}

func validateMiddleName(middleName string) error {
	if middleName == "" {
		return nil
	}

	if err := fnsfio.ValidateMiddleName(middleName); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	return nil
}
