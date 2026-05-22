package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("attempt usecase forbidden")
	ErrInvalidInput = errors.New("attempt usecase invalid input")
	ErrLimitReached = errors.New("attempt usecase limit reached")
)

func RequireStudent(actor role.Role) error {
	if actor.Kind() != role.TypeStudent {
		return fmt.Errorf("%w: only student can work with quiz attempts", ErrForbidden)
	}
	return nil
}
