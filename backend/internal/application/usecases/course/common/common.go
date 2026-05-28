package common

import (
	"context"
	"errors"
	"fmt"
	"strings"

	enrollmentports "gitflic.ru/lms/backend/internal/application/ports/enrollment"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

var (
	ErrForbidden    = errors.New("доступ к курсу запрещён")
	ErrInvalidInput = errors.New("некорректные параметры запроса к курсу")
)

func RequireManager(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeCreator:
		return nil
	default:
		return fmt.Errorf("%w: управление курсами доступно только администратору и создателю", ErrForbidden)
	}
}

func RequireStudent(actor role.Role) error {
	if actor.Kind() != role.TypeStudent {
		return fmt.Errorf("%w: просмотр записанных курсов доступен только студенту", ErrForbidden)
	}
	return nil
}

func RequireProgressAccess(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeOrganization, role.TypeStudent:
		return nil
	default:
		return fmt.Errorf("%w: доступ к данным прогресса разрешён только администратору, организации и студенту", ErrForbidden)
	}
}

func RequireOrganization(actor role.Role) error {
	switch actor.Kind() {
	case role.TypeAdmin, role.TypeOrganization:
		return nil
	default:
		return fmt.Errorf("%w: доступ к ресурсу разрешён только администратору и организации", ErrForbidden)
	}
}

func RequireOrganizationScope(ctx context.Context, scope enrollmentports.OrganizationScope, actor role.Role, actorPersonID, enrollmentID string) error {
	if actor.Kind() != role.TypeOrganization {
		return nil
	}
	personID, err := uuid.Parse(actorPersonID)
	if err != nil {
		return fmt.Errorf("%w: некорректный идентификатор пользователя при проверке области организации", ErrForbidden)
	}
	enrID, err := uuid.Parse(enrollmentID)
	if err != nil {
		return fmt.Errorf("%w: некорректный идентификатор зачисления при проверке области организации", ErrForbidden)
	}
	belongs, err := scope.EnrollmentBelongsToPersonOrganization(ctx, enrID, personID)
	if err != nil {
		return fmt.Errorf("verify organization scope: %w", err)
	}
	if !belongs {
		return fmt.Errorf("%w: организация может просматривать прогресс только своих студентов", ErrForbidden)
	}
	return nil
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

func NormalizeSearchText(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}
