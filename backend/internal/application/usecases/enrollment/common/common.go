package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("enrollment usecase forbidden")
	ErrInvalidInput = errors.New("enrollment usecase invalid input")
	ErrConflict     = errors.New("enrollment usecase conflict")
)

func RequireManager(actor role.Role) error {
	if actor.Kind() != role.TypeAdmin {
		return fmt.Errorf("%w: only admin can manage enrollments", ErrForbidden)
	}
	return nil
}
