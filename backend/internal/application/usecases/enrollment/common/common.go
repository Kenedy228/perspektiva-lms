package common

import (
	"errors"
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/role"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
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

func RequireViewer(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeStudent:
		return nil
	default:
		return fmt.Errorf("%w: просмотр зачислений доступен только администратору и студенту", ErrForbidden)
	}
}

func NormalizePagination(limit, offset int) (int, int, error) {
	if offset < 0 {
		return 0, 0, fmt.Errorf("%w: смещение не может быть отрицательным", ErrInvalidInput)
	}
	if limit == 0 {
		return DefaultLimit, offset, nil
	}
	if limit < 0 {
		return 0, 0, fmt.Errorf("%w: лимит не может быть отрицательным", ErrInvalidInput)
	}
	if limit > MaxLimit {
		return 0, 0, fmt.Errorf("%w: лимит не может превышать %d", ErrInvalidInput, MaxLimit)
	}
	return limit, offset, nil
}
