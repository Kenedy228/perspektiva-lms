package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("question usecase forbidden")
	ErrInvalidInput = errors.New("question usecase invalid input")
)

func RequireAuthor(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: only admin or creator can manage questions", ErrForbidden)
	}
}
