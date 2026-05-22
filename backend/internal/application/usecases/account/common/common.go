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
	ErrForbidden          = errors.New("account usecase forbidden")
	ErrInvalidInput       = errors.New("account usecase invalid input")
	ErrConflict           = errors.New("account usecase conflict")
	ErrInvalidCredentials = errors.New("invalid account credentials")
)

func RequireAdmin(actor role.Role) error {
	if actor.Kind() != role.TypeAdmin {
		return fmt.Errorf("%w: only admin can change accounts", ErrForbidden)
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
