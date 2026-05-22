package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("quiz usecase forbidden")
	ErrInvalidInput = errors.New("quiz usecase invalid input")
)

func RequireManager(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: only admin or creator can manage quizzes", ErrForbidden)
	}
}
