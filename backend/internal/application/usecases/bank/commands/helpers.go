package commands

import (
	"context"
	"fmt"

	bankports "gitflic.ru/lms/backend/internal/application/ports/bank"
	"gitflic.ru/lms/backend/internal/application/usecases/bank/common"
	bankdomain "gitflic.ru/lms/backend/internal/domain/bank"
	"gitflic.ru/lms/backend/internal/domain/bank/title"
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

func parseRequiredUUIDs(values []string, field string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(values))
	for i := range values {
		id, err := parseRequiredUUID(values[i], field)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func buildTitle(value string) (title.Title, error) {
	t, err := title.New(value)
	if err != nil {
		return title.Title{}, fmt.Errorf("create bank title: %w", err)
	}
	return t, nil
}

func loadBank(ctx context.Context, r bankports.Repository, id string) (*bankdomain.Bank, error) {
	bankID, err := parseRequiredUUID(id, "bank id")
	if err != nil {
		return nil, err
	}

	b, err := r.FindByID(ctx, bankID)
	if err != nil {
		return nil, fmt.Errorf("find bank: %w", err)
	}

	return b, nil
}

func saveBank(ctx context.Context, r bankports.Repository, b *bankdomain.Bank) error {
	if err := r.Save(ctx, b); err != nil {
		return fmt.Errorf("save bank: %w", err)
	}
	return nil
}

func recordAudit(ctx context.Context, audit bankports.AuditRecorder, action bankports.AuditAction, bankID string, actor role.Role) error {
	if audit == nil {
		return nil
	}

	if err := audit.RecordBankAudit(ctx, bankports.AuditEvent{
		Action:    action,
		BankID:    bankID,
		ActorRole: actor.Kind().String(),
	}); err != nil {
		return fmt.Errorf("record bank audit: %w", err)
	}

	return nil
}
