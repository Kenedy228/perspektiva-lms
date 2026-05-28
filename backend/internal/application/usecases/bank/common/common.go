package common

import (
	"errors"
	"fmt"
	"strings"

	"gitflic.ru/lms/backend/internal/domain/role"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

var (
	ErrForbidden    = errors.New("доступ к банку вопросов запрещён")
	ErrInvalidInput = errors.New("некорректные параметры запроса к банку вопросов")
)

func RequireManager(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: управлять банками вопросов могут только администратор и создатель", ErrForbidden)
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
		return 0, 0, fmt.Errorf("%w: размер страницы не может быть отрицательным", ErrInvalidInput)
	}

	if limit > MaxLimit {
		return 0, 0, fmt.Errorf("%w: размер страницы не должен превышать %d", ErrInvalidInput, MaxLimit)
	}

	return limit, offset, nil
}

func NormalizeSearchText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}
