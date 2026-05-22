package commands

import (
	"context"
	"fmt"

	accountports "gitflic.ru/lms/backend/internal/application/ports/account"
	"gitflic.ru/lms/backend/internal/application/usecases/account/common"
	accountdomain "gitflic.ru/lms/backend/internal/domain/account"
	"gitflic.ru/lms/backend/internal/domain/account/login"
	"gitflic.ru/lms/backend/internal/domain/role"
	"github.com/google/uuid"
)

func parseRequiredUUID(value, field string) (uuid.UUID, error) {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse %s: %w", field, err)
	}

	if id == uuid.Nil {
		return uuid.Nil, fmt.Errorf("%w: %s is required", common.ErrInvalidInput, field)
	}

	return id, nil
}

func requireLoginAvailable(ctx context.Context, r accountports.Repository, l login.Login, exclude uuid.UUID) error {
	exists, err := r.LoginExists(ctx, l, exclude)
	if err != nil {
		return fmt.Errorf("check account login uniqueness: %w", err)
	}

	if exists {
		return fmt.Errorf("%w: login already belongs to another account", common.ErrConflict)
	}

	return nil
}

func requirePersonWithoutAccount(ctx context.Context, r accountports.Repository, personID uuid.UUID) error {
	exists, err := r.PersonHasAccount(ctx, personID, uuid.Nil)
	if err != nil {
		return fmt.Errorf("check person account uniqueness: %w", err)
	}

	if exists {
		return fmt.Errorf("%w: person already has an account", common.ErrConflict)
	}

	return nil
}

func recordAudit(ctx context.Context, audit accountports.AuditRecorder, action accountports.AuditAction, acc *accountdomain.Account, actor role.Role) error {
	if err := audit.RecordAccountAudit(ctx, accountports.AuditEvent{
		Action:    action,
		AccountID: acc.ID().String(),
		PersonID:  acc.PersonID().String(),
		ActorRole: actor.Kind().String(),
	}); err != nil {
		return fmt.Errorf("record account audit: %w", err)
	}

	return nil
}
