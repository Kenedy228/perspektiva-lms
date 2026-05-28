package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("доступ к тесту запрещён")
	ErrInvalidInput = errors.New("некорректные параметры запроса к тесту")
)

func RequireManager(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: управлять тестами могут только администратор и создатель", ErrForbidden)
	}
}
