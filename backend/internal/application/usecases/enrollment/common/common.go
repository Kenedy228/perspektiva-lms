package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("доступ к зачислению запрещён")
	ErrInvalidInput = errors.New("некорректные параметры запроса к зачислению")
	ErrConflict     = errors.New("конфликт при зачислении")
)

func RequireManager(actor role.Role) error {
	if actor.Kind() != role.TypeAdmin {
		return fmt.Errorf("%w: управление зачислениями доступно только администратору", ErrForbidden)
	}
	return nil
}
