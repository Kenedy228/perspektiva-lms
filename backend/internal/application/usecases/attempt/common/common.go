package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

var (
	ErrForbidden    = errors.New("доступ к попытке запрещён")
	ErrInvalidInput = errors.New("некорректные параметры запроса к попытке")
	ErrLimitReached = errors.New("превышен лимит попыток")
)

func RequireStudent(actor role.Role) error {
	if actor.Kind() != role.TypeStudent {
		return fmt.Errorf("%w: работать с попытками теста может только студент", ErrForbidden)
	}
	return nil
}
