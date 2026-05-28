package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("доступ к вопросу запрещён")
	ErrInvalidInput = errors.New("некорректные параметры запроса к вопросу")
)

func RequireAuthor(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: управлять вопросами могут только администратор и создатель", ErrForbidden)
	}
}
