package bank

import (
	"context"

	bankdomain "gitflic.ru/lms/backend/internal/domain/bank"
	"github.com/google/uuid"
)

type Repository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*bankdomain.Bank, error)
	Save(ctx context.Context, b *bankdomain.Bank) error
	DeleteByID(ctx context.Context, id uuid.UUID) error
}

type AuditAction string

const (
	AuditActionCreate          AuditAction = "bank.create"
	AuditActionRename          AuditAction = "bank.rename"
	AuditActionAddQuestions    AuditAction = "bank.add_questions"
	AuditActionRemoveQuestions AuditAction = "bank.remove_questions"
	AuditActionClearQuestions  AuditAction = "bank.clear_questions"
	AuditActionDelete          AuditAction = "bank.delete"
)

type AuditEvent struct {
	Action    AuditAction
	BankID    string
	ActorRole string
}

type AuditRecorder interface {
	RecordBankAudit(ctx context.Context, event AuditEvent) error
}
