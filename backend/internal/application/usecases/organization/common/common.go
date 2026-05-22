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
	ErrForbidden    = errors.New("organization usecase forbidden")
	ErrInvalidInput = errors.New("organization usecase invalid input")
)

func RequireAdmin(actor role.Role) error {
	if actor.Kind() != role.TypeAdmin {
		return fmt.Errorf("%w: only admin can interact with organization usecases", ErrForbidden)
	}

	return nil
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

func NormalizeINNSearch(value string) string {
	return strings.ReplaceAll(strings.TrimSpace(value), " ", "")
}
