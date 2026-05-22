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
	ErrForbidden    = errors.New("person usecase forbidden")
	ErrInvalidInput = errors.New("person usecase invalid input")
	ErrConflict     = errors.New("person usecase conflict")
)

func RequireAdmin(actor role.Role) error {
	if actor.Kind() != role.TypeAdmin {
		return fmt.Errorf("%w: only admin can interact with person usecases", ErrForbidden)
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

func NormalizeSNILSSearch(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "-", "")
	return value
}
