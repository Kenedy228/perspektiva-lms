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
	ErrForbidden    = errors.New("bank usecase forbidden")
	ErrInvalidInput = errors.New("bank usecase invalid input")
)

func RequireManager(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: only admin or creator can manage question banks", ErrForbidden)
	}
}

func NormalizePagination(limit, offset int) (int, int, error) {
	if offset < 0 {
		return 0, 0, fmt.Errorf("%w: offset cannot be negative", ErrInvalidInput)
	}

	if limit == 0 {
		return DefaultLimit, offset, nil
	}

	if limit < 0 {
		return 0, 0, fmt.Errorf("%w: limit cannot be negative", ErrInvalidInput)
	}

	if limit > MaxLimit {
		return 0, 0, fmt.Errorf("%w: limit cannot exceed %d", ErrInvalidInput, MaxLimit)
	}

	return limit, offset, nil
}

func NormalizeSearchText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}
