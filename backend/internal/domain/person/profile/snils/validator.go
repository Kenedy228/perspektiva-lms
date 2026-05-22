package snils

import (
	"fmt"

	snils2 "github.com/Kenedy228/snils"
)

func validate(value string) error {
	if err := snils2.Validate(value); err != nil {
		return fmt.Errorf("%w: %w", ErrInvalid, err)
	}

	return nil
}
