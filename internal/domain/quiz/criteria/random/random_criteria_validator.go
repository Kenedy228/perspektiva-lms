package random

import "errors"

var (
	ErrInvalidCount = errors.New("invalid criteria count")
)

func validateCount(count int) error {
	if count <= 0 {
		return ErrInvalidCount
	}

	return nil
}

